// Package daemon implements the core engine for Revere's daemon mode.
package daemon

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/yext/revere/db"
	"github.com/yext/revere/env"
)

// Daemon represents the part of Revere that actually executes monitors and
// triggers and dispatches alerts.
type Daemon struct {
	monitors map[db.MonitorID]*monitor

	lastMonitorsUpdate time.Time

	stop    chan struct{}
	stopper sync.Once
	stopped chan struct{}

	*env.Env
}

// New initializes a new Daemon. To actually make the Daemon run, call Start.
func New(env *env.Env) *Daemon {
	return &Daemon{
		monitors: make(map[db.MonitorID]*monitor),
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
		Env:      env,
	}
}

// Start starts running a Daemon.
func (d *Daemon) Start() {
	go d.run()
}

func (d *Daemon) run() {
	defer close(d.stopped)

	log.Info("Daemon is running.")

	t := time.NewTicker(time.Duration(10) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			d.updateMonitors()
		case <-d.stop:
			return
		}
	}
}

func (d *Daemon) updateMonitors() {
	currentUpdateTime := time.Now()

	threshold := d.lastMonitorsUpdate
	if !threshold.IsZero() {
		// Provide some buffer for clock skew.
		threshold = threshold.Add(time.Duration(-5) * time.Minute)
	}

	infos, err := d.DB.LoadMonitorVersionInfosUpdatedSince(threshold)
	if err != nil {
		log.WithError(err).Error("Could not load list of updated monitors.")
	}

	for _, info := range infos {
		old := d.monitors[info.MonitorID]
		if old != nil {
			if old.version == info.Version && info.Archived == nil {
				// Already running newest version.
				continue
			}

			log.WithFields(log.Fields{
				"monitor": old.id,
				"version": old.version,
			}).Info("Tearing down monitor.")

			old.stop()
			delete(d.monitors, info.MonitorID)
		}

		if info.Archived != nil {
			// Don't run archived monitors.
			continue
		}

		log.WithFields(log.Fields{
			"monitor": info.MonitorID,
			"version": info.Version,
		}).Info("Starting monitor.")

		new, err := newMonitor(info.MonitorID, d.Env)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"monitor": info.MonitorID,
			}).Error("Load monitor failed.")

			// TODO(eefi): Put in a placeholder that constantly
			// reports _:Unknown to Revere admin. For now, just
			// ignore this monitor.
			continue
		}

		d.monitors[info.MonitorID] = new
		new.start()
	}

	d.lastMonitorsUpdate = currentUpdateTime
}

// Stop gracefully stops a Daemon. It tries to allow any in-progress delivery of
// alerts to finish before returning.
func (d *Daemon) Stop() {
	d.stopper.Do(func() {
		// Stop run loop first to avoid race over what monitors exist to
		// be stopped.
		close(d.stop)
		<-d.stopped

		for id, m := range d.monitors {
			log.WithFields(log.Fields{
				"monitor": m.id,
				"version": m.version,
			}).Info("Tearing down monitor.")

			m.stop()
			delete(d.monitors, id)
		}

		log.Info("Daemon has stopped.")
	})
}

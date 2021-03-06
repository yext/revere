package probe

import (
	"sync"
	"time"

	"github.com/juju/errors"
)

// Polling helps implement probes that check conditions at regular intervals.
// Embed a pointer to this struct in a struct that implements Checker to get
// such a polling probe without having to deal with the polling loop.
type Polling struct {
	period       time.Duration
	checker      Checker
	readingsSink chan<- []Reading

	stop    chan struct{}
	stopper sync.Once
	stopped chan struct{}
}

func NewPolling(period time.Duration, checker Checker, readingsSink chan<- []Reading) (*Polling, error) {
	if period <= 0 {
		return nil, errors.Errorf("cannot poll with nonpositive period %s", period)
	}

	return &Polling{
		period:       period,
		checker:      checker,
		readingsSink: readingsSink,
		stop:         make(chan struct{}),
		stopped:      make(chan struct{}),
	}, nil
}

func (p *Polling) Start() {
	go p.poll()
}

func (p *Polling) Stop() {
	p.stopper.Do(func() {
		close(p.stop)
		<-p.stopped
	})
}

func (p *Polling) poll() {
	defer close(p.stopped)

	t := time.NewTicker(p.period)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			p.readingsSink <- p.checker.Check()
		case <-p.stop:
			return
		}
	}
}

// Checker is used by Polling to actually check the thing being monitored.
type Checker interface {
	Check() []Reading
}

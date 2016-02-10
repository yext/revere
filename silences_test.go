package revere_test

import (
	"testing"
	"time"

	. "github.com/yext/revere"
)

var (
	silenceEndLimit = 14 * 24 * time.Hour
)

func presentSilence() *Silence {
	now := time.Now()
	s := new(Silence)
	s.Id = 1
	s.Subprobes = "test.*.example"
	s.MonitorId = 1
	s.MonitorName = "test monitor"
	s.Start = now.AddDate(0, 0, -1)
	s.End = now.AddDate(0, 0, 1)
	return s
}

func pastSilence() *Silence {
	s := presentSilence()
	s.Start = s.Start.AddDate(0, 0, -1)
	s.End = s.End.AddDate(0, 0, -2)
	return s
}

func futureSilence() *Silence {
	s := presentSilence()
	s.Start = s.Start.AddDate(0, 0, 2)
	s.End = s.End.AddDate(0, 0, 1)
	return s
}

func TestIsPastSilence(t *testing.T) {
	s := pastSilence()
	if s.IsPast(time.Now()) == false {
		t.Errorf("Expected past silence for silence starting at %v and ending at %v\n", s.Start, s.End)
	}
}

func TestIsPresentSilence(t *testing.T) {
	s := presentSilence()
	if s.IsPresent(time.Now()) == false {
		t.Errorf("Expected present silence for silence starting at %v and ending at %v\n", s.Start, s.End)
	}
}

func TestEditableSilence(t *testing.T) {
	s := futureSilence()
	if s.Editable() == false {
		t.Errorf("Expected editable silence for silence starting at %v and ending at %v\n", s.Start, s.End)
	}
}

func TestEditingPastSilence(t *testing.T) {
	s := futureSilence()
	ps := pastSilence()
	errs := s.ValidateAgainstOld(ps)
	if errs == nil {
		t.Error("Expected error trying to edit past silence")
	}
}

func TestPastSilenceCreate(t *testing.T) {
	ps := pastSilence()
	errs := ps.ValidateAgainstOld(nil)
	if errs == nil {
		t.Error("Expected error trying to create a silence in the past")
	}
}

func TestCreateSilenceInvalidMonitorId(t *testing.T) {
	fs := futureSilence()
	fs.MonitorId = 0
	errs := fs.ValidateAgainstOld(nil)
	if errs == nil {
		t.Error("Expected error trying to create a silence with an invalid monitor id")
	}
}

func TestEditSilenceInvalidMonitorId(t *testing.T) {
	s := presentSilence()
	s.MonitorId = 0
	errs := s.ValidateAgainstOld(futureSilence())
	if errs == nil {
		t.Error("Expected error trying to edit a silence with a different monitor id")
	}
}

func TestEditSilenceInvalidSubprobes(t *testing.T) {
	s := presentSilence()
	s.Subprobes = ""
	errs := s.ValidateAgainstOld(futureSilence())
	if errs == nil {
		t.Error("Expected error trying to edit a silence with a different subprobe")
	}
}

func TestCreateSilenceStartAfterEnd(t *testing.T) {
	fs := futureSilence()
	fs.Start = fs.End.AddDate(0, 0, 1)
	errs := fs.ValidateAgainstOld(nil)
	if errs == nil {
		t.Error("Expected error trying to create a silence with a start after the end")
	}
}

func TestSilenceInvalidDuration(t *testing.T) {
	fs := futureSilence()
	fs.End = fs.End.Add(silenceEndLimit)
	errs := fs.ValidateAgainstOld(nil)
	if errs == nil {
		t.Error("Expected error trying to create a silence with a end date beyond the allowed limit")
	}
}

func TestPresentSilenceStartEdit(t *testing.T) {
	fs := futureSilence()
	errs := fs.ValidateAgainstOld(presentSilence())
	if errs == nil {
		t.Error("Expected error trying to set a start time for a currently running silence")
	}
}

func TestValidSilenceCreate(t *testing.T) {
	fs := futureSilence()
	errs := fs.ValidateAgainstOld(nil)
	if errs != nil {
		t.Errorf("Unexpected error trying to create a new silence: %v", errs)
	}
}

func TestValidPresentSilenceEdit(t *testing.T) {
	oldPs := presentSilence()
	newPs := presentSilence()
	newPs.Start = oldPs.Start
	newPs.End = newPs.End.AddDate(0, 0, 1)
	errs := newPs.ValidateAgainstOld(oldPs)
	if errs != nil {
		t.Errorf("Unexpected error trying to edit a present silence: %v", errs)
	}
}

func TestValidFutureSilenceEdit(t *testing.T) {
	oldFs := futureSilence()
	newFs := futureSilence()
	newFs.Start = newFs.Start.AddDate(0, 0, 1)
	newFs.End = newFs.End.AddDate(0, 0, 1)
	errs := newFs.ValidateAgainstOld(oldFs)
	if errs != nil {
		t.Errorf("Unexpected error trying to edit a future silence: %v", errs)
	}
}

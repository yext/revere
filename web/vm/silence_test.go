package vm

import (
	"testing"
	"time"
)

func presentSilence() *Silence {
	now := time.Now()
	s := new(Silence)
	s.SilenceId = 1
	s.Subprobe = "test.*.example"
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
	old := pastSilence()
	s := presentSilence()
	errs := append(s.validate(), s.validateOld(old)...)
	if errs == nil {
		t.Error("Expected error trying to edit past silence")
	}
}

func TestPastSilenceCreate(t *testing.T) {
	s := pastSilence()
	errs := append(s.validate(), s.validateNew()...)
	if errs == nil {
		t.Error("Expected error trying to create a silence in the past")
	}
}

func TestCreateSilenceInvalidMonitorId(t *testing.T) {
	s := futureSilence()
	s.MonitorId = 0
	errs := append(s.validate(), s.validateNew()...)
	if errs == nil {
		t.Error("Expected error trying to create a silence with an invalid monitor id")
	}
}

func TestEditSilenceInvalidMonitorId(t *testing.T) {
	old := futureSilence()
	s := presentSilence()
	s.MonitorId = 2
	errs := append(s.validate(), s.validateOld(old)...)
	if errs == nil {
		t.Error("Expected error trying to edit a silence with a different monitor id")
	}
}

func TestEditSilenceInvalidSubprobes(t *testing.T) {
	old := presentSilence()
	s := presentSilence()
	s.Subprobe = ""
	errs := append(s.validate(), s.validateOld(old)...)
	if errs == nil {
		t.Error("Expected error trying to edit a silence with a different subprobe")
	}
}

func TestCreateSilenceStartAfterEnd(t *testing.T) {
	s := futureSilence()
	s.Start = s.End.AddDate(0, 0, 1)
	errs := append(s.validate(), s.validateNew()...)
	if errs == nil {
		t.Error("Expected error trying to create a silence with a start after the end")
	}
}

func TestSilenceInvalidDuration(t *testing.T) {
	s := futureSilence()
	s.End = s.End.Add(maxSilenceDuration)
	errs := append(s.validate(), s.validateNew()...)
	if errs == nil {
		t.Error("Expected error trying to create a silence with a end date beyond the allowed limit")
	}
}

func TestPresentSilenceStartEdit(t *testing.T) {
	old := presentSilence()
	s := futureSilence()
	errs := append(s.validate(), s.validateOld(old)...)
	if errs == nil {
		t.Error("Expected error trying to set a start time for a currently running silence")
	}
}

func TestValidSilenceCreate(t *testing.T) {
	s := futureSilence()
	errs := append(s.validate(), s.validateNew()...)
	if errs != nil {
		t.Errorf("Unexpected error trying to create a new silence: %v", errs)
	}
}

func TestValidPresentSilenceEdit(t *testing.T) {
	old := presentSilence()
	s := presentSilence()
	s.Start = old.Start
	s.End = s.End.AddDate(0, 0, 1)
	errs := append(s.validate(), s.validateOld(old)...)
	if errs != nil {
		t.Errorf("Unexpected error trying to edit a present silence: %v", errs)
	}
}

func TestValidFutureSilenceEdit(t *testing.T) {
	old := futureSilence()
	s := futureSilence()
	s.Start = s.Start.AddDate(0, 0, 1)
	s.End = s.End.AddDate(0, 0, 1)
	errs := append(s.validate(), s.validateOld(old)...)
	if errs != nil {
		t.Errorf("Unexpected error trying to edit a future silence: %v", errs)
	}
}

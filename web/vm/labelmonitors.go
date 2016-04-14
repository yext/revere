package vm

import "github.com/yext/revere"

type LabelMonitor struct {
	*revere.LabelMonitor
}

func (lm *LabelMonitor) Id() int64 {
	return int64(lm.LabelMonitor.MonitorId)
}

func (lm *LabelMonitor) Name() string {
	return lm.LabelMonitor.Name
}

func (lm *LabelMonitor) Description() string {
	return lm.LabelMonitor.Description
}

func (lm *LabelMonitor) Subprobe() string {
	return lm.LabelMonitor.Subprobe
}

func NewLabelMonitors(labelMonitors []*revere.LabelMonitor) []*LabelMonitor {
	viewmodels := make([]*LabelMonitor, len(labelMonitors))
	for i, lm := range labelMonitors {
		viewmodels[i] = &LabelMonitor{lm}
	}
	return viewmodels
}

func BlankLabelMonitors() []*LabelMonitor {
	return []*LabelMonitor{}
}

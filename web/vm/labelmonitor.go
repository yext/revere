package vm

import "github.com/yext/revere"

type LabelMonitor struct {
	*revere.LabelMonitor
}

func (lm *LabelMonitor) Id() int64 {
	return int64(lm.LabelMonitor.MonitorId)
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
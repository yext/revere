package vm

import (
	"fmt"
)

type Breadcrumb struct {
	Text string
	Link string
}

func MonitorIndexBcs() []Breadcrumb {
	return []Breadcrumb{Breadcrumb{"Monitors", "/monitors"}}
}

func MonitorViewBcs(mn string, id int64) []Breadcrumb {
	return append(MonitorIndexBcs(), Breadcrumb{mn, fmt.Sprintf("/monitors/%d", id)})
}

func SubprobeIndexBcs(mn string, id int64) []Breadcrumb {
	return append(MonitorViewBcs(mn, id), Breadcrumb{"Subprobe", fmt.Sprintf("/monitors/%d/subprobes", id)})
}

func SubprobeViewBcs(s *Subprobe) []Breadcrumb {
	return append(SubprobeIndexBcs(s.MonitorName, int64(s.MonitorID)), Breadcrumb{s.Name, fmt.Sprintf("/monitors/%d/subprobes/%d", s.MonitorID, s.SubprobeID)})
}

func SilencesIndexBcs() []Breadcrumb {
	return []Breadcrumb{Breadcrumb{"Silences", "/silences"}}
}

func SilencesViewBcs(id int64, mn string) []Breadcrumb {
	return append(SilencesIndexBcs(), Breadcrumb{fmt.Sprintf("Silence for %s", mn), fmt.Sprintf("/silences/%d", id)})
}

func LabelIndexBcs() []Breadcrumb {
	return []Breadcrumb{Breadcrumb{"Labels", "/labels"}}
}

func LabelViewBcs(mn string, id int64) []Breadcrumb {
	return append(LabelIndexBcs(), Breadcrumb{mn, fmt.Sprintf("/labels/%d", id)})
}

func IsLastBc(a []Breadcrumb, i int) bool {
	return i == len(a)-1
}

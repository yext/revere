package web

import (
	"fmt"

	"github.com/yext/revere"
)

type breadcrumb struct {
	Text string
	Link string
}

func monitorIndexBcs() []breadcrumb {
	return []breadcrumb{breadcrumb{"monitors", "/monitors"}}
}

func monitorViewBcs(mn string, id uint) []breadcrumb {
	return append(monitorIndexBcs(), breadcrumb{mn, fmt.Sprintf("/monitors/%d", id)})
}

func subprobeIndexBcs(mn string, id uint) []breadcrumb {
	return append(monitorViewBcs(mn, id), breadcrumb{"subprobe", fmt.Sprintf("/monitors/%d/subprobes", id)})
}

func subprobeViewBcs(s *revere.Subprobe) []breadcrumb {
	return append(subprobeIndexBcs(s.MonitorName, s.MonitorId), breadcrumb{s.Name, fmt.Sprintf("/monitors/%d/subprobes/%d", s.MonitorId, s.Id)})
}

func silencesIndexBcs() []breadcrumb {
	return []breadcrumb{breadcrumb{"silences", "/silences"}}
}

func silencesViewBcs(id uint, mn string) []breadcrumb {
	return append(silencesIndexBcs(), breadcrumb{fmt.Sprintf("silence for %s", mn), fmt.Sprintf("/silences/%d", id)})
}

func isLastBc(a []breadcrumb, i int) bool {
	return i == len(a)-1
}

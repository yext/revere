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

func monitorViewBcs(m *revere.Monitor) []breadcrumb {
	return append(monitorIndexBcs(), breadcrumb{m.Name, fmt.Sprintf("/monitors/%d", m.Id)})
}

func subprobeIndexBcs(m *revere.Monitor) []breadcrumb {
	return append(monitorViewBcs(m), breadcrumb{"subprobe", fmt.Sprintf("/monitors/%d/subprobes", m.Id)})
}

func subprobeViewBcs(m *revere.Monitor, s *revere.Subprobe) []breadcrumb {
	return append(subprobeIndexBcs(m), breadcrumb{s.Name, fmt.Sprintf("/monitors/%d/subprobes/%d", m.Id, s.Id)})
}

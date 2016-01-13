package web

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	probesScriptDir   = "web/js/probes/"
	probesServingPath = "/static/js/probes/"

	targetsScriptDir   = "web/js/targets/"
	targetsServingPath = "/static/js/targets/"
)

func getScripts(dir string, servingPath string) (scripts []string, err error) {
	scriptInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, s := range scriptInfo {
		if s.IsDir() || !strings.HasSuffix(s.Name(), ".js") {
			continue
		}
		scripts = append(scripts, fmt.Sprintf("%s%s", servingPath, s.Name()))
	}
	return scripts, nil
}

func probeScripts() (scripts []string, err error) {
	return getScripts(probesScriptDir, probesServingPath)
}

func targetScripts() (scripts []string, err error) {
	return getScripts(targetsScriptDir, targetsServingPath)
}

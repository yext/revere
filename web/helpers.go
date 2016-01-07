package web

import (
	"fmt"
	"io/ioutil"
)

const (
	targetsScriptDir   = "web/js/targets/"
	targetsServingPath = "/static/js/targets/"
)

func getScripts(dir string, servingPath string) (scripts []string, err error) {
	scriptInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, s := range scriptInfo {
		if s.IsDir() {
			continue
		}
		scripts = append(scripts, fmt.Sprintf("%s%s", servingPath, s.Name()))
	}
	return scripts, nil
}

func targetScripts() (scripts []string, err error) {
	return getScripts(targetsScriptDir, targetsServingPath)
}

package vm

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
)

const (
	defaultScriptTag = "<script type=\"text/javascript\" src=\"/%s\"></script>"

	baseDir         = "web/js"
	baseServingPath = "static/js"

	probesDir         = "web/js/probes"
	probesServingPath = "static/js/probes"

	targetsDir         = "web/js/targets"
	targetsServingPath = "static/js/targets"
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

func ProbeScripts() (scripts []string, err error) {
	return getScripts(probesDir, probesServingPath)
}

func TargetScripts() (scripts []string, err error) {
	return getScripts(targetsDir, targetsServingPath)
}

func newScripts(scripts []string) template.HTML {
	return addScripts(scripts, template.HTML(""))
}

func addScripts(addl []string, current template.HTML) template.HTML {
	buffer := bytes.Buffer{}
	for _, script := range addl {
		buffer.WriteString(fmt.Sprintf(defaultScriptTag, path.Join(baseServingPath, script)))
		buffer.WriteString("\n")
	}
	buffer.WriteString(string(current))
	return template.HTML(buffer.String())
}

func newScript(script string) template.HTML {
	return addScript(script, template.HTML(""))
}

func addScript(addl string, current template.HTML) template.HTML {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf(defaultScriptTag, path.Join(baseServingPath, addl)))
	buffer.WriteString("\n")
	buffer.WriteString(string(current))
	return template.HTML(buffer.String())
}

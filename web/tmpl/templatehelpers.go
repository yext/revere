// Package tmpl provides a wrapper around golang's html/template package, and
// provides additional template manipulation functions.
package tmpl

import (
	"html/template"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/juju/errors"
	"github.com/yext/revere/boxes"
)

type Template struct {
	htmlTemplate *template.Template
}

var (
	partialsStrings map[string]string      = make(map[string]string)
	functions       map[string]interface{} = make(map[string]interface{})

	baseServingPath = "static/js" // URL
	partials        = "partials"

	htmlExt = regexp.MustCompile(`^(?i).+\.html$`)
)

func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.htmlTemplate.Execute(w, data)
}

func (t *Template) AddTmpl(name string) {
	tmplBox := boxes.Templates()
	tmplString := tmplBox.MustString(name)

	_ = template.Must(t.htmlTemplate.New(path.Base(name)).Parse(tmplString))
}

func (t *Template) AddTmpls(names []string) {
	for _, name := range names {
		t.AddTmpl(name)
	}
}

func AddDefaultFunc(name string, function interface{}) {
	functions[name] = function
}

func SetPartialsLocation(location string) {
	partials = location
}

func NewTemplate(name string) *Template {
	// put in pull request https://github.com/GeertJohan/go.rice/pull/82
	// eventually would like to change this to MustFindBox
	tmplBox := boxes.Templates()
	tmplString := tmplBox.MustString(name)
	t := template.Must(template.New(path.Base(name)).Funcs(functions).Parse(tmplString))

	// So that partials is only read once
	if len(partialsStrings) == 0 {
		populatePartials()
	}

	for name, str := range partialsStrings {
		_ = template.Must(t.New(name).Funcs(functions).Parse(str))
	}

	return &Template{t}
}

// This is called because go.rice states that FindBox cannot be called in an
// init method
func populatePartials() {
	tmplBox := boxes.Templates()
	err := tmplBox.Walk(partials, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Trace(err)
		}

		if !htmlExt.MatchString(filepath) {
			return nil
		}

		if !info.IsDir() {
			partialsStrings[path.Base(filepath)] = tmplBox.MustString(filepath)
		}

		return nil
	})
	if err != nil {
		panic(errors.Trace(err))
	}
}

func SetTitle(data map[string]interface{}, title string) map[string]interface{} {
	data["Title"] = title
	return data
}

func StrEq(a, b interface{}) bool {
	c, ok := a.(string)
	if !ok {
		return false
	}
	d, ok := b.(string)
	if !ok {
		return false
	}

	return c == d
}

func GetScript(filepath string) string {
	return path.Join(baseServingPath, filepath)
}

func AppendDir(dir string, scripts []string) []string {
	result := make([]string, len(scripts))
	for i, script := range scripts {
		result[i] = path.Join(dir, script)
	}
	return result
}

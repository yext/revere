package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
)

type Template struct {
	htmlTemplate *template.Template
}

var (
	partials  string
	functions map[string]interface{} = make(map[string]interface{})
)

func (t Template) Name() string {
	return t.htmlTemplate.Name()
}

func (t Template) Execute(w io.Writer, data interface{}) error {
	return t.htmlTemplate.Execute(w, data)
}

func (t Template) GetTmpl() *template.Template {
	return t.htmlTemplate
}

func AddDefaultFunc(name string, function interface{}) {
	functions[name] = function
}

func SetPartialsLocation(location string) {
	partials = location
}

func InitTemplates(dir string, funcs template.FuncMap) (templates map[string]*template.Template) {
	templates = make(map[string]*template.Template)
	templateInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("Got error initializing templates: %v", err))
	}
	for _, t := range templateInfo {
		if t.IsDir() || !strings.HasSuffix(t.Name(), ".html") {
			continue
		}
		templates[t.Name()] = template.Must(template.New(t.Name()).Funcs(funcs).ParseFiles(dir + t.Name()))
	}
	return
}

func NewTemplate(dir string, name string) *Template {
	t := template.Must(template.New(name).Funcs(functions).ParseFiles(dir + name))
	t = template.Must(t.ParseGlob(partials))

	return &Template{t}
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
	if c == d {
		return true
	}
	return false
}

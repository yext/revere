package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Template struct {
	htmlTemplate *template.Template
}

var (
	partials  string
	functions map[string]interface{} = make(map[string]interface{})
)

const (
	templatesDir string = "web/views/"
)

func (t *Template) Name() string {
	return t.htmlTemplate.Name()
}

func (t *Template) Execute(w io.Writer, data interface{}) error {
	return t.htmlTemplate.Execute(w, data)
}

func (t *Template) GetTmpl() *template.Template {
	return t.htmlTemplate
}

func (t *Template) AddFunc(name string, function interface{}) {
	t.htmlTemplate.Funcs(map[string]interface{}{name: function})
}

func (t *Template) AddTmpl(name string) {
	t.htmlTemplate = template.Must(t.htmlTemplate.ParseFiles(path.Join(templatesDir, name)))
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

func InitTemplates(dir string, funcs template.FuncMap) (templates map[string]*template.Template) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Cannot get the rooted path name of the current directory"))
	}

	// When running packages for testing
	if _, err = ioutil.ReadDir(path.Join(pwd, "web")); err != nil {
		dir = fmt.Sprintf("../%s", dir)
	}

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

func NewTemplateDir(dir string, name string) *Template {
	return newTemplate(path.Join(templatesDir, dir, name))
}

func NewTemplate(name string) *Template {
	return newTemplate(path.Join(templatesDir, name))
}

func newTemplate(filepath string) *Template {
	t := template.Must(template.New(path.Base(filepath)).Funcs(functions).ParseFiles(filepath))
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

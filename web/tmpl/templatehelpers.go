package tmpl

import (
	"html/template"
	"io"
	"path"
)

type Template struct {
	htmlTemplate *template.Template
}

var (
	partials  string
	functions map[string]interface{} = make(map[string]interface{})
)

const (
	baseServingPath        = "static/js"
	baseDir                = "web/js"
	templatesDir    string = "web/views/"
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
	if c == d {
		return true
	}
	return false
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

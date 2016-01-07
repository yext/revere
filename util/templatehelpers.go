package util

import (
	"fmt"
	"html/template"
	"io/ioutil"
)

func InitTemplates(dir string, funcs template.FuncMap) (tMap map[string]*template.Template) {
	tMap = make(map[string]*template.Template)
	templateInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("Got error initializing templates: %v", err))
	}
	for _, t := range templateInfo {
		if t.IsDir() {
			continue
		}
		tMap[t.Name()], err = template.New(t.Name()).Funcs(funcs).ParseFiles(dir + t.Name())
		if err != nil {
			panic(fmt.Sprintf("Got error initializing templates: %v", err))
		}
	}
	return
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

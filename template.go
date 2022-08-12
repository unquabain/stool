package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/masterminds/sprig"
	yaml "gopkg.in/yaml.v3"
)

func FuncMap() template.FuncMap {
	fm := sprig.TxtFuncMap()
	fm[`yaml`] = func(v any) (string, error) {
		b, e := yaml.Marshal(v)
		return string(b), e
	}
	fm[`yml`] = fm[`yaml`]
	fm[`json`] = func(v any) (string, error) {
		b, e := json.Marshal(v)
		return string(b), e
	}
	fm[`js`] = fm[`json`]
	fm[`jsonpretty`] = func(v any) (string, error) {
		b, e := json.MarshalIndent(v, ``, `    `)
		return string(b), e
	}
	fm[`jspretty`] = fm[`jsonpretty`]
	return fm
}

func GetTemplate(ttext, tfile string) (*template.Template, error) {
	t := template.New(`cmdline`)
	t = t.Funcs(FuncMap())
	if tfile != `` {
		tfilebytes, err := os.ReadFile(tfile)
		if err != nil {
			return nil, fmt.Errorf(`could not read template file %q: %w`, tfile, err)
		}
		ttext = string(tfilebytes)
	}
	t, err := t.Parse(`{{ range . }}` + ttext + `{{ end }}`)
	if err != nil {
		return nil, fmt.Errorf(`unable to parse template: %w`, err)
	}
	return t, nil
}

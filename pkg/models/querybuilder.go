package models

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

var fns = template.FuncMap{
	"last": func(x int, a interface{}) bool {
		return x == reflect.ValueOf(a).Len()-1
	},
	"inc": func(i int) int {
		return i + 1
	},
	"vault": func(v string) string {
		return "!vault |\n" + v
	},
	"ident": func(spaces int, v string) string {
		pad := strings.Repeat(" ", spaces)
		return pad + strings.Replace(v, "\n", "\n"+pad, -1)
	},
}

type QueryBuilder struct {
	Sort      string
	Order     string
	Range     []int64
	Page      int64
	RefKey    string
	Id        []int64
	FilterKey string
	FilterVal string
	Columns   []string
}

func (obj QueryBuilder) List(query string) (*bytes.Buffer, error) {
	sql := query + `
	{{ if gt (len .Id) 0 }}
	where
		{{ .RefKey }} in (?)
	{{ else if gt (len .FilterVal) 0 }}
	where
		{{ .FilterKey }} = "{{ .FilterVal }}"
	{{ end }}
	{{ if gt (len .Sort) 0 }}
	order by {{ .Sort }} {{ .Order }}
	{{ end }}
	{{ if eq (len .Range) 2 }}
	limit {{ index .Range 0 }}, {{ .Page }}
	{{ end }}`

	var tpl bytes.Buffer
	t := template.Must(template.New("list").Parse(sql))
	err := t.Execute(&tpl, obj)
	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

func (obj QueryBuilder) Insert(table string) (*bytes.Buffer, error) {
	sql := `
	insert into ` + table + `(
		{{ range $i, $key := . }}
			{{ $key }}{{if last $i $}}{{else}},{{end}}{{ end }}
	) values (
		{{ range $i, $key := . }}
			:{{ $key }}{{if last $i $}}{{else}},{{end}}{{ end }}
	)`

	var tpl bytes.Buffer
	t := template.Must(template.New("insert").Funcs(fns).Parse(sql))
	err := t.Execute(&tpl, obj.Columns)
	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

func (obj QueryBuilder) Update(table string) (*bytes.Buffer, error) {
	sql := `
	update
		` + table + `
	set
	{{ range $i, $key := . }}
		{{ $key }} = :{{ $key }}{{if last $i $}}{{else}},{{end}}{{ end }}
	where
		id = :id
	`

	var tpl bytes.Buffer
	t := template.Must(template.New("update").Funcs(fns).Parse(sql))
	err := t.Execute(&tpl, obj.Columns)
	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

func (obj QueryBuilder) Delete(table string) (*bytes.Buffer, error) {
	sql := `
	delete from
		` + table + `
	where
		id = :id
	`

	var tpl bytes.Buffer
	t := template.Must(template.New("update").Parse(sql))
	err := t.Execute(&tpl, obj.Columns)
	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

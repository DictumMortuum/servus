package models

import (
	"bytes"
	"text/template"
)

type Args struct {
	Sort      string
	Order     string
	Range     []int64
	Page      int64
	RefKey    string
	Id        []int64
	FilterKey string
	FilterVal string
}

func (obj Args) List(query string) (*bytes.Buffer, error) {
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

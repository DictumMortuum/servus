package models

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"reflect"
	"strconv"
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
	Context   *gin.Context
	Sort      string
	Order     string
	Range     []int64
	Page      int64
	RefKey    string
	Id        int64
	Ids       []int64
	FilterKey string
	FilterVal string
	Columns   []string
	Data      map[string]interface{}
	Resources map[string]bool
}

func NewArgsFromContext(c *gin.Context) (*QueryBuilder, error) {
	qb := QueryBuilder{}
	return qb.NewFromContext(c)
}

func (obj QueryBuilder) NewFromContext(c *gin.Context) (*QueryBuilder, error) {
	rs := QueryBuilder{
		Page:      5,
		Resources: map[string]bool{},
		Columns:   []string{},
		Context:   c,
	}

	arg := c.Params.ByName("id")

	if arg != "" {
		id, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return nil, err
		}

		rs.Id = id
	}

	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		c.BindJSON(&rs.Data)

		for key := range rs.Data {
			rs.Columns = append(rs.Columns, key)
		}
	}

	args := c.Request.URL.Query()

	if val, ok := args["_sort"]; ok {
		rs.Sort = val[0]
	}

	if val, ok := args["_order"]; ok {
		rs.Order = val[0]
	} else {
		rs.Order = "ASC"
	}

	if val, ok := args["_start"]; ok {
		start, err := strconv.ParseInt(val[0], 10, 64)
		if err == nil {
			rs.Range = append(rs.Range, start)
		}
	}

	if val, ok := args["_end"]; ok {
		end, err := strconv.ParseInt(val[0], 10, 64)
		if err == nil {
			rs.Range = append(rs.Range, end)
		}
	}

	if len(rs.Range) == 2 {
		rs.Page = rs.Range[1] - rs.Range[0] + 1
	}

	if val, ok := args["id"]; ok {
		for _, raw := range val {
			id, err := strconv.ParseInt(raw, 10, 64)
			if err == nil {
				rs.Ids = append(rs.Ids, id)
			}
			rs.RefKey = "id"
		}
	}

	if val, ok := args["_resource"]; ok {
		for _, raw := range val {
			rs.Resources[raw] = true
		}
	}

	for key, val := range args {
		if strings.HasSuffix(key, "id") && key != "id" {
			// this is a reference
			id, err := strconv.ParseInt(val[0], 10, 64)
			if err == nil {
				rs.RefKey = key
				rs.Ids = append(rs.Ids, id)
			}
		} else if !strings.HasPrefix(key, "_") && key != "id" {
			// this is a filter
			rs.FilterKey = key
			rs.FilterVal = val[0]
		}
	}

	return &rs, nil
}

func (obj QueryBuilder) List(query string) (*bytes.Buffer, error) {
	sql := query + `
	{{ if gt (len .Ids) 0 }}
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

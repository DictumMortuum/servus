package boardgames

import (
	"bytes"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/DictumMortuum/servus/pkg/boardgames/deta"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"regexp"
	"strings"
	"text/template"
)

var (
	ignore_brackets = regexp.MustCompile(`(?s)\((.*)\)`)
	gamescom        = regexp.MustCompile(`\S+ Επιτραπέζιο (.*) για`)
	epitrapezio     = regexp.MustCompile(`[^|]+| [0-9]+ Ετών`)
)

func TransformName(s string) string {
	pre := gamescom.FindStringSubmatch(s)
	var tmp string

	if len(pre) == 2 {
		tmp = pre[1]
	} else {
		tmp = s
	}

	pre2 := epitrapezio.FindStringSubmatch(tmp)

	if len(pre2) == 1 {
		tmp = pre2[0]
	}

	tmp = ignore_brackets.ReplaceAllString(tmp, "")
	tmp = strings.ToLower(tmp)
	tmp = strings.TrimSpace(tmp)
	tmp = strings.ReplaceAll(tmp, "-", " ")
	tmp = strings.ReplaceAll(tmp, ",", " ")
	tmp = strings.ReplaceAll(tmp, "&", " ")
	tmp = strings.ReplaceAll(tmp, "–", " ")
	tmp = strings.ReplaceAll(tmp, ":", " ")
	tmp = strings.ReplaceAll(tmp, "/", " ")
	tmp = strings.ReplaceAll(tmp, "expansion", " ")
	tmp = strings.ReplaceAll(tmp, "Expansion", " ")
	tmp = strings.ReplaceAll(tmp, " KS ", " ")
	tmp = strings.ReplaceAll(tmp, " ks ", " ")
	tmp = strings.ReplaceAll(tmp, "Kickstarter", " ")
	tmp = strings.ReplaceAll(tmp, "kickstarter", " ")
	tmp = strings.ReplaceAll(tmp, " exp ", " ")
	tmp = strings.ReplaceAll(tmp, " exp. ", " ")
	tmp = strings.ReplaceAll(tmp, "επιτραπέζιο παιχνίδι", "")
	tmp = strings.ReplaceAll(tmp, "στρατηγικής", "")
	tmp = strings.ReplaceAll(tmp, "σε μεταλλικό κουτί", "")
	tmp = strings.ReplaceAll(tmp, "κλασικό παιχνίδι", "")
	tmp = strings.ReplaceAll(tmp, "επέκταση για", "")
	tmp = strings.ReplaceAll(tmp, "οικογενειακό", "")
	tmp = strings.ReplaceAll(tmp, "παιχνίδι", "")
	tmp = strings.ReplaceAll(tmp, "με τράπουλα", "")
	tmp = strings.ReplaceAll(tmp, "συνεργατικό", "")
	tmp = strings.ReplaceAll(tmp, "στρατηγικό", "")
	tmp = strings.ReplaceAll(tmp, "παράρτημα για", "")
	tmp = strings.ReplaceAll(tmp, "παράρτημα επιτραπέζιου παιχνιδιού", "")
	tmp = strings.ReplaceAll(tmp, "επιτραπέζιο", "")
	tmp = strings.ReplaceAll(tmp, "κάισσα", "")
	tmp = strings.ReplaceAll(tmp, "lcg", "")
	tmp = strings.ReplaceAll(tmp, "επέκταση", "")
	tmp = strings.ReplaceAll(tmp, " en ", "")
	return tmp
}

func GetPriceById(db *sqlx.DB, id int64) (*models.Price, error) {
	var rs models.Price

	err := db.QueryRowx(`
		select
			p.*,
			g.rank,
			IFNULL(g.thumb,"") as thumb,
			IFNULL(g.preview,"") as preview,
			IFNULL(g.name,"") as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.id = ? and
			p.store_id = s.id
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	rs.Key = fmt.Sprintf("%d", rs.Id)
	rs.TransformedName = TransformName(rs.Name)

	return &rs, nil
}

func GetHistoricPricesById(db *sqlx.DB, id int64) ([]models.HistoricPrice, error) {
	var rs []models.HistoricPrice

	sql := `
		select
			boardgame_id,
			cr_date,
			avg(price) avg,
			min(price) min,
			max(price) max
		from
			tboardgamepriceshistory
		where
			boardgame_id = ?
		group by 1,YEAR(cr_date),MONTH(cr_date)
	`

	err := db.Select(&rs, sql, id)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func GetPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	return GetPriceById(db, args.Id)
}

func countPrices(db *sqlx.DB, args *models.QueryBuilder) (int, error) {
	var tpl bytes.Buffer

	sql := `
		select
			1
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.store_id = s.id and
			p.ignored = 0
		{{ if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = {{ .FilterVal }}
		{{ end }}
	`

	t := template.Must(template.New("count").Parse(sql))
	err := t.Execute(&tpl, args)
	if err != nil {
		return -1, err
	}

	var count []int
	err = db.Select(&count, tpl.String())
	if err != nil {
		return -1, err
	}

	return len(count), nil
}

func GetListPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs []models.Price

	count, err := countPrices(db, args)
	if err != nil {
		return nil, err
	}
	args.Context.Header("X-Total-Count", fmt.Sprintf("%d", count))

	sql := `
		select
			p.*,
			g.rank,
			IFNULL(g.thumb,"") as thumb,
			IFNULL(g.preview,"") as preview,
			IFNULL(g.name,"") as boardgame_name,
			s.name as store_name
		from
			tboardgameprices p
			left join tboardgames g on g.id = p.boardgame_id,
			tboardgamestores s
		where
			p.store_id = s.id and
			p.ignored = 0
		{{ if gt (len .Ids) 0 }}
			and p.{{ .RefKey }} in (?)
		{{ else if gt (len .FilterVal) 0 }}
			and p.{{ .FilterKey }} = {{ .FilterVal }}
		{{ end }}
		{{ if gt (len .Sort) 0 }}
		order by {{ .Sort }} {{ .Order }}
		{{ else }}
		order by g.rank asc, p.cr_date asc
		{{ end }}
		{{ if eq (len .Range) 2 }}
		limit {{ index .Range 0 }}, {{ .Page }}
		{{ end }}`

	var tpl bytes.Buffer
	t := template.Must(template.New("list").Parse(sql))
	err = t.Execute(&tpl, args)
	if err != nil {
		return nil, err
	}

	query, ids, err := sqlx.In(tpl.String(), args.Ids)
	if err != nil {
		query = tpl.String()
	} else {
		query = db.Rebind(query)
	}

	err = db.Select(&rs, query, ids...)
	if err != nil {
		return nil, err
	}

	if args.Resources["deta"] {
		deta_db, err := deta.New("prices")
		if err != nil {
			return nil, err
		}

		for _, item := range rs {
			item.Key = fmt.Sprintf("%d", item.Id)
			err = deta.Populate(deta_db, item)
			if err != nil {
				return nil, err
			}
		}
	}

	if args.Resources["history"] {
		retval := []models.Price{}

		for _, item := range rs {
			item.HistoricPrices, err = GetHistoricPricesById(db, item.BoardgameId.Int64)
			if err != nil {
				return nil, err
			}
			retval = append(retval, item)
		}
		return retval, nil
	} else {
		retval := []models.Price{}

		for _, item := range rs {
			item.TransformedName = TransformName(item.Name)
			retval = append(retval, item)
		}

		return retval, nil
	}
}

func CreatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	var rs models.Price

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		err := fn(args.Data, true)
		if err != nil {
			return nil, err
		}
	}

	query, err := args.Insert("tboardgameprices")
	if err != nil {
		return nil, err
	}

	price, err := db.NamedExec(query.String(), &rs)
	if err != nil {
		return nil, err
	}

	id, err := price.LastInsertId()
	if err != nil {
		return nil, err
	}

	rs.Id = id
	return rs, nil
}

func UpdatePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	args.IgnoreColumn("rank")
	args.IgnoreColumn("thumb")
	args.IgnoreColumn("boardgame_name")
	args.IgnoreColumn("store_name")
	args.IgnoreColumn("transformed_name")
	args.IgnoreColumn("key")
	args.IgnoreColumn("preview")

	updateFns := rs.Constructor()
	for _, fn := range updateFns {
		fn(args.Data, false)
	}

	if rs.BoardgameId.Valid {
		exists, err := boardgameExists(db, map[string]interface{}{
			"id": rs.BoardgameId.Int64,
		})
		if err != nil {
			return nil, err
		}

		if exists == nil {
			_, err := bgg.FetchBoardgame(db, rs.BoardgameId.Int64)
			if err != nil {
				return nil, err
			}
		}
	}

	sql, err := args.Update("tboardgameprices")
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(sql.String(), &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func DeletePrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`delete from tboardgameprices where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func UnmapPrice(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs, err := GetPriceById(db, args.Id)
	if err != nil {
		return nil, err
	}

	_, err = db.NamedExec(`update tboardgameprices set boardgame_id = NULL where id = :id`, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

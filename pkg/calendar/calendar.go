package calendar

import (
	"errors"
	"github.com/DictumMortuum/servus/pkg/db"
	"github.com/tealeg/xlsx/v3"
	"regexp"
	"strconv"
)

type XlsxHeaderInfo struct {
	Index   int
	Name    string
	Info    int
	Enabled bool
}

type XlsxHeader struct {
	Id         string `xlsx:"0"`
	Surname    string `xlsx:"1"`
	Name       string `xlsx:"2"`
	Speciality string `xlsx:"3"`
	RowIndex   int
	Info       []XlsxHeaderInfo
}

type Person struct {
	Id         string `xlsx:"0"`
	Surname    string `xlsx:"1"`
	Name       string `xlsx:"2"`
	Speciality string `xlsx:"3"`
	Calendar   []db.CalendarRow
}

type Sheet struct {
	Header XlsxHeader
	People []Person
}

func (c *Person) New(row *xlsx.Row, header XlsxHeader) error {
	row.ReadStruct(c)
	skips := 4

	if c.Id == "" || c.Name == "" || c.Surname == "" || c.Speciality == "" {
		return errors.New("not a person entry")
	}

	for j := skips; j < row.Sheet.MaxCol; j++ {
		cell := row.GetCell(j)
		raw, _ := cell.FormattedValue()
		i := row.GetCoordinate()

		if header.Info[j-4].Enabled {
			c.Calendar = append(c.Calendar, db.CalendarRow{
				Index:   header.Info[j-4].Index,
				DayName: header.Info[j-4].Name,
				Raw:     raw,
				Row:     row,
				X:       i,
				Y:       j,
			})
		} else {
			skips++
		}
	}

	return nil
}

func GetCoworkers(d db.CalendarRow, header XlsxHeader) []string {
	coworkers := []string{}
	sheet := d.Row.Sheet

	for i := header.RowIndex + 1; i < sheet.MaxRow; i++ {
		temp, _ := sheet.Cell(i, d.Y)
		t, _ := temp.FormattedValue()

		if i != d.X {
			if d.Raw == t && RawToShift(d.Raw) >= 0 {
				var person Person
				r, _ := sheet.Row(i)
				r.ReadStruct(&person)
				coworkers = append(coworkers, person.Name+" "+person.Surname)
			}
		}
	}

	return coworkers
}

func (h *XlsxHeader) New(row *xlsx.Row) error {
	row.ReadStruct(h)

	if h.Id != "ΜΗΤΡΩΑ" && h.Surname != "ΕΠΙΘΕΤΟ" && h.Name != "ΟΝΟΜΑ" && h.Speciality != "ΕΙΔΙΚΟΤΗΤΑ" {
		return errors.New("not header")
	}

	h.RowIndex = row.GetCoordinate()

	for j := 4; j < row.Sheet.MaxCol; j++ {
		cell, _ := row.GetCell(j).FormattedValue()
		re := regexp.MustCompile("([0-9]+)([ΔΕΤΡPΤTΕΠΕEΠΑΣΑAΚKΥY]{2})([0-9]*)")
		s := re.FindStringSubmatch(cell)

		if len(s) > 0 {
			s1, _ := strconv.Atoi(s[1])
			s3, _ := strconv.Atoi(s[3])
			h.Info = append(h.Info, XlsxHeaderInfo{s1, s[2], s3, true})
		} else {
			h.Info = append(h.Info, XlsxHeaderInfo{-1, "", -1, false})
		}
	}

	return nil
}

func (s *Sheet) New(file string) error {
	wb, err := xlsx.OpenFile(file)
	if err != nil {
		return err
	}

	sheet := wb.Sheets[0]
	header := XlsxHeader{}

	for i := 0; i < sheet.MaxRow; i++ {
		row, err := sheet.Row(i)
		if err != nil {
			return err
		}

		err = header.New(row)
		if err != nil {
			continue
		} else {
			s.Header = header
			break
		}
	}

	for i := header.RowIndex + 1; i < sheet.MaxRow; i++ {
		current, _ := sheet.Row(i)
		person := Person{}

		err = person.New(current, header)
		if err != nil {
			continue
		} else {
			s.People = append(s.People, person)
		}
	}

	return nil
}

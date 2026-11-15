package seating

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/importandrew/MaCRa/models"
)

// csvFields are the columns of a seating file: tour, table, east, south,
// west, north.
const csvFields = 2 + models.SeatsPerTable

// ParseCSV reads a seating prepared for several tours: one row per table,
// "tour,table,east,south,west,north", the seats as participant numbers.
// The separator may be a comma, a semicolon (as spreadsheets save CSV in
// many locales) or a tab, whichever the first line has; a first line that
// doesn't start with a number is a header. Tables of a tour go 1, 2, 3, …
// and a number sits once per tour. Whether the numbers match the
// participants is checked when a tour is seated (Fixed).
func ParseCSV(text string) (models.PlanTours, error) {
	text = strings.TrimPrefix(text, "\uFEFF")
	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = separator(text)
	reader.FieldsPerRecord = -1 // checked below, with the line
	reader.TrimLeadingSpace = true

	file := planFile{plan: models.PlanTours{}, seated: map[int]map[int]bool{}}
	for first := true; ; first = false {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line, _ := reader.FieldPos(0)
		if err != nil {
			if parseErr, ok := errors.AsType[*csv.ParseError](err); ok {
				line = parseErr.Line
			}
			return nil, models.WithData(models.ErrSeatingFileRow, models.ErrorData{dataLine: line})
		}
		if isBlank(record) || first && !startsWithNumber(record) {
			continue // a blank row, or a header
		}
		if err := file.add(record, line); err != nil {
			return nil, err
		}
	}
	return file.finish()
}

// planFile is a seating file read so far: the plan, and per tour the
// numbers seated.
type planFile struct {
	plan   models.PlanTours
	seated map[int]map[int]bool
}

// add adds a row (a table) of the file at the line.
func (f *planFile) add(record []string, line int) error {
	values, ok := rowNumbers(record)
	if !ok {
		return models.WithData(models.ErrSeatingFileRow, models.ErrorData{dataLine: line})
	}
	tour, table := values[0], values[1]
	tables := f.plan[tour]
	for len(tables) < table {
		tables = append(tables, [4]int{})
	}
	if tables[table-1] != [4]int{} {
		return models.WithData(models.ErrSeatingFileTableTwice,
			models.ErrorData{dataLine: line, dataTour: tour, "table": table})
	}
	if f.seated[tour] == nil {
		f.seated[tour] = map[int]bool{}
	}
	for _, number := range values[2:] {
		if f.seated[tour][number] {
			return models.WithData(models.ErrSeatingFileNumberTwice,
				models.ErrorData{dataLine: line, dataTour: tour, "number": number})
		}
		f.seated[tour][number] = true
	}
	tables[table-1] = [4]int(values[2:])
	f.plan[tour] = tables
	return nil
}

// finish returns the plan, refusing an empty one or a tour that skips a
// table number.
func (f *planFile) finish() (models.PlanTours, error) {
	if len(f.plan) == 0 {
		return nil, models.ErrSeatingFileEmpty
	}
	for tour, tables := range f.plan {
		for i, table := range tables {
			if table == [4]int{} {
				return nil, models.WithData(models.ErrSeatingFileTableMissing,
					models.ErrorData{dataTour: tour, "table": i + 1})
			}
		}
	}
	return f.plan, nil
}

// maxFileTable bounds a table number in a seating file, so a typo can't
// make a huge plan.
const maxFileTable = 1000

// rowNumbers returns a row's six values, each a whole number from 1 (a
// table up to maxFileTable), or false.
func rowNumbers(record []string) ([]int, bool) {
	if len(record) < csvFields {
		return nil, false
	}
	for _, extra := range record[csvFields:] {
		if strings.TrimSpace(extra) != "" {
			return nil, false
		}
	}
	values := make([]int, csvFields)
	for i, field := range record[:csvFields] {
		n, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil || n < 1 || n > maxFileTable*models.SeatsPerTable {
			return nil, false
		}
		values[i] = n
	}
	if values[1] > maxFileTable {
		return nil, false
	}
	return values, true
}

// separator returns the first line's separator: a semicolon or a tab if it
// has one and no comma, else a comma.
func separator(text string) rune {
	firstLine, _, _ := strings.Cut(text, "\n")
	if !strings.Contains(firstLine, ",") {
		if strings.Contains(firstLine, ";") {
			return ';'
		}
		if strings.Contains(firstLine, "\t") {
			return '\t'
		}
	}
	return ','
}

// isBlank reports whether a row has no values (e.g. ";;;;;" from a
// spreadsheet).
func isBlank(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

// startsWithNumber reports whether a row's first value is a number.
func startsWithNumber(record []string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(record[0]))
	return err == nil
}

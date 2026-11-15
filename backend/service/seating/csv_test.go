package seating

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/models"
)

func TestParseCSV(t *testing.T) {
	want := models.PlanTours{
		1: {{1, 2, 3, 4}, {5, 6, 7, 8}},
		2: {{1, 5, 3, 7}, {2, 6, 4, 8}},
	}
	for name, file := range map[string]string{
		"comma":     "1,1,1,2,3,4\n1,2,5,6,7,8\n2,1,1,5,3,7\n2,2,2,6,4,8\n",
		"header":    "tour,table,east,south,west,north\n1,1,1,2,3,4\n1,2,5,6,7,8\n2,1,1,5,3,7\n2,2,2,6,4,8",
		"semicolon": "\uFEFFТур;Стол;В;Ю;З;С\r\n2;2;2;6;4;8\r\n1;2;5;6;7;8\r\n;;;;;\r\n1;1;1;2;3;4\r\n2;1;1;5;3;7\r\n",
		"tab":       "1\t1\t1\t2\t3\t4\n1\t2\t5\t6\t7\t8\n\n2\t1\t 1\t5\t3\t7\n2\t2\t2\t6\t4\t8\t\n",
	} {
		t.Run(name, func(t *testing.T) {
			plan, err := ParseCSV(file)
			require.NoError(t, err)
			assert.Equal(t, want, plan)
		})
	}
}

func TestParseCSV_errors(t *testing.T) {
	for name, c := range map[string]struct {
		file string
		err  error
		data models.ErrorData
	}{
		"empty":         {"", models.ErrSeatingFileEmpty, nil},
		"header only":   {"tour,table,east,south,west,north\n", models.ErrSeatingFileEmpty, nil},
		"five values":   {"1,1,1,2,3,4\n1,2,5,6,7\n", models.ErrSeatingFileRow, models.ErrorData{"line": 2}},
		"seven values":  {"1,1,1,2,3,4,9\n", models.ErrSeatingFileRow, models.ErrorData{"line": 1}},
		"a name":        {"1,1,1,2,3,4\n1,2,5,6,7,Bob\n", models.ErrSeatingFileRow, models.ErrorData{"line": 2}},
		"zero":          {"1,1,0,2,3,4\n", models.ErrSeatingFileRow, models.ErrorData{"line": 1}},
		"second header": {"1,1,1,2,3,4\ntour,table,e,s,w,n\n", models.ErrSeatingFileRow, models.ErrorData{"line": 2}},
		"bad quote":     {"1,1,1,2,3,4\n1,\"2,5,6,7,8\n", models.ErrSeatingFileRow, models.ErrorData{"line": 2}},
		"table twice": {"1,1,1,2,3,4\n1,1,5,6,7,8\n", models.ErrSeatingFileTableTwice,
			models.ErrorData{"line": 2, "tour": 1, "table": 1}},
		"number twice": {"1,1,1,2,3,4\n1,2,5,6,7,1\n", models.ErrSeatingFileNumberTwice,
			models.ErrorData{"line": 2, "tour": 1, "number": 1}},
		"table missing": {"1,1,1,2,3,4\n2,1,1,2,3,4\n2,3,5,6,7,8\n", models.ErrSeatingFileTableMissing,
			models.ErrorData{"tour": 2, "table": 2}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ParseCSV(c.file)
			require.ErrorIs(t, err, c.err)
			if c.data != nil {
				assertData(t, err, c.data)
			}
		})
	}
}

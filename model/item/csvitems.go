package item

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// csvItems is a struct that helps read raw item data in a map-based format.
// The first line of the raw data must be a column name map. This line is not
// accessible to the user. Subsequent lines will be available through the
// get(string) function.
type csvItems struct {
	*csv.Reader
	columnMap map[string]int
	record    []string
	err       error
}

// get returns the value in the column with the label k.
func (c *csvItems) get(k string) string {
	i, ok := c.columnMap[k]
	if !ok {
		panic(fmt.Errorf("tried to fetch nonexistent column %s from %v", k, *c))
	}
	return c.record[i]
}

func (c *csvItems) next() bool {
	c.record, c.err = c.Read()
	return c.err == nil
}

func newCsvMapReader(reader *csv.Reader) (*csvItems, error) {
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}

	columns := make(map[string]int)
	for i, colname := range record {
		columns[colname] = i
	}
	return &csvItems{Reader: reader, columnMap: columns}, nil
}

func mustMapCsv(s string) *csvItems {
	csvReader := csv.NewReader(strings.NewReader(s))
	csvReader.Comment = '#'
	r, err := newCsvMapReader(csvReader)
	if err != nil {
		panic(err)
	}
	return r
}

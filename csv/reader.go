package csv

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Rows is the records of a csv file. Its cursor starts before
// the first row of the result set. Use Next to advance from row to row.
type Rows struct {
	fields   []string
	records  [][]string
	length   int
	position int
	closed   bool
}

// ReadAll reads all the remaining records from r.
func ReadAll(r io.Reader) (*Rows, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}

	rs := &Rows{}
	switch len(records) {
	case 0:
		return nil, fmt.Errorf("Empty csv file")
	case 1:
		rs.fields = records[0]
		return rs, nil
	default:
		rs.fields = records[0]
		rs.records = records[1:]
		rs.length = len(rs.records)
	}

	return rs, nil
}

// Fields returns the fieldnames.
func (rs *Rows) Fields() []string {
	return rs.fields
}

// Next prepares the next result row for reading with the Scan method.
func (rs *Rows) Next() bool {
	if rs.position < rs.length {
		rs.position++
		return true
	}
	rs.closed = true
	return false
}

// Scan copies the columns in the current row into the values pointed at by dest.
// The number of values in dest must be the same as the number of columns in Rows.
func (rs *Rows) Scan(dest ...*string) error {
	if rs.closed {
		return fmt.Errorf("Rows are closed")
	}
	if len(dest) != len(rs.fields) {
		return fmt.Errorf("expected %d destination arguments in Scan, not %d", len(rs.fields), len(dest))
	}

	for i, v := range rs.records[rs.position-1] {
		*dest[i] = v
	}

	return nil
}

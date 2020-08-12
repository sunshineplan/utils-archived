package export

import (
	"encoding/csv"
	"fmt"
	"io"
)

func exportCSV(fieldnames []string, rows []interface{}, w io.Writer, utf8bom bool) error {
	if utf8bom {
		w.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	writer := csv.NewWriter(w)
	if err := writer.Write(fieldnames); err != nil {
		return err
	}
	for _, row := range rows {
		m, ok := row.(map[string]interface{})
		if !ok {
			return fmt.Errorf("can not convert row to map[string]interface{}")
		}
		r := make([]string, len(fieldnames))
		for index, fieldname := range fieldnames {
			if m[fieldname] != nil {
				r[index] = fmt.Sprintf("%v", m[fieldname])
			}
		}
		if err := writer.Write(r); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

// CSV export csv content to writer
func CSV(fieldnames []string, rows []interface{}, w io.Writer) error {
	return exportCSV(fieldnames, rows, w, false)
}

// CSVWithUTF8BOM export csv content to writer
func CSVWithUTF8BOM(fieldnames []string, rows []interface{}, w io.Writer) error {
	return exportCSV(fieldnames, rows, w, true)
}

package export

import (
	"encoding/csv"
	"fmt"
	"io"
)

// CSV export csv content to writer
func CSV(fieldnames []string, rows []interface{}, writer io.Writer) error {
	w := csv.NewWriter(writer)
	if err := w.Write(fieldnames); err != nil {
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
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
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
		r := make([]string, len(fieldnames))
		switch obj := reflect.ValueOf(row); obj.Kind() {
		case reflect.Map:
			if reflect.TypeOf(row).Key().Name() == "string" {
				for index, fieldname := range fieldnames {
					if v := obj.MapIndex(reflect.ValueOf(fieldname)); v.IsValid() {
						r[index] = fmt.Sprintf("%v", v)
					}
				}
			} else {
				return fmt.Errorf("can not export rows which map is not string")
			}
		case reflect.Struct:
			for index, fieldname := range fieldnames {
				if v := obj.FieldByName(fieldname); v.IsValid() {
					r[index] = fmt.Sprintf("%v", v)
				}
			}
		default:
			return fmt.Errorf("not support rows format: %s", obj.Kind())
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

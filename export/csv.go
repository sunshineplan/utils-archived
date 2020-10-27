package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func exportCSV(fieldnames []string, slice interface{}, w io.Writer, utf8bom bool) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return fmt.Errorf("rows is not slice")
	}
	rows := reflect.ValueOf(slice)
	if fieldnames == nil {
		var err error
		fieldnames, err = getStructFieldNames(rows.Index(0).Interface())
		if err != nil {
			return err
		}
	}

	if utf8bom {
		w.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	writer := csv.NewWriter(w)
	if err := writer.Write(fieldnames); err != nil {
		return err
	}

	for i := 0; i < rows.Len(); i++ {
		row := rows.Index(i)
		if row.Kind() == reflect.Interface {
			row = row.Elem()
		}
		r := make([]string, len(fieldnames))
		switch row.Kind() {
		case reflect.Map:
			if reflect.TypeOf(row.Interface()).Key().Name() == "string" {
				for index, fieldname := range fieldnames {
					if v := row.MapIndex(reflect.ValueOf(fieldname)); v.IsValid() && v.Interface() != nil {
						r[index] = fmt.Sprintf("%v", v)
					}
				}
			} else {
				return fmt.Errorf("can not export rows which map is not string")
			}
		case reflect.Struct:
			for index, fieldname := range fieldnames {
				if v := row.FieldByName(fieldname); v.IsValid() && v.Interface() != nil {
					r[index] = fmt.Sprintf("%v", v)
				}
			}
		default:
			return fmt.Errorf("not support rows format: %s", row.Kind())
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
func CSV(fieldnames []string, slice interface{}, w io.Writer) error {
	return exportCSV(fieldnames, slice, w, false)
}

// CSVWithUTF8BOM export csv content to writer
func CSVWithUTF8BOM(fieldnames []string, slice interface{}, w io.Writer) error {
	return exportCSV(fieldnames, slice, w, true)
}

func getStructFieldNames(i interface{}) ([]string, error) {
	if reflect.TypeOf(i).Kind() != reflect.Struct {
		return nil, fmt.Errorf("can not get fieldnames from interface which is not struct")
	}
	v := reflect.ValueOf(i)
	var fieldnames []string
	for i := 0; i < v.NumField(); i++ {
		fieldnames = append(fieldnames, v.Type().Field(i).Name)
	}
	return fieldnames, nil
}

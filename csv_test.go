package utils

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCSV(t *testing.T) {
	type test struct{ A, B interface{} }
	testcase := []struct {
		name      string
		filenames []string
		slice     interface{}
	}{
		{
			name:      "map slice",
			filenames: []string{"A", "B"},
			slice: []map[string]interface{}{
				map[string]interface{}{"A": "a", "B": "b"},
				map[string]interface{}{"A": "aa", "B": nil},
			},
		},
		{
			name:      "struct slice",
			filenames: []string{"A", "B"},
			slice:     []test{test{A: "a", B: "b"}, test{A: "aa", B: nil}},
		},
		{
			name:      "struct slice without filenames",
			filenames: nil,
			slice:     []test{test{A: "a", B: "b"}, test{A: "aa", B: nil}},
		},
		{
			name:      "interface slice",
			filenames: []string{"A", "B"},
			slice: []interface{}{
				test{A: "a", B: "b"},
				map[string]interface{}{"A": "aa", "B": nil},
			},
		},
	}
	result := `A,B
a,b
aa,
`

	for _, tc := range testcase {
		var b bytes.Buffer
		if err := ExportCSV(tc.filenames, tc.slice, &b); err != nil {
			t.Error(tc.name, err)
		}
		if r := b.String(); r != result {
			t.Errorf("%s expected %q; got %q", tc.name, result, r)
		}
	}
}

func TestCSVWithUTF8BOM(t *testing.T) {
	result := `A,B
a,b
`
	var b bytes.Buffer
	if err := ExportUTF8CSV([]string{"A", "B"}, []interface{}{map[string]string{"A": "a", "B": "b"}}, &b); err != nil {
		t.Error(err)
	}
	c := b.Bytes()
	if !reflect.DeepEqual(utf8bom, c[:3]) {
		t.Errorf("expected %q; got %q", utf8bom, c[:3])
	}
	if r := string(c[3:]); r != result {
		t.Errorf("expected %q; got %q", result, r)
	}
}

func TestGetStructFieldNames(t *testing.T) {
	if _, err := getStructFieldNames(map[string]string{"test": "test"}); err == nil {
		t.Error("gave nil error; want error")
	}
	if fieldnames, err := getStructFieldNames(struct{ A, B string }{}); err != nil {
		t.Error(err)
	} else {
		if !reflect.DeepEqual([]string{"A", "B"}, fieldnames) {
			t.Errorf("expected %q; got %q", []string{"A", "B"}, fieldnames)
		}
	}
}

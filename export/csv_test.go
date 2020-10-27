package export

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestCSV(t *testing.T) {
	type test struct{ A, B interface{} }
	mapSlice := []map[string]interface{}{map[string]interface{}{"A": "a", "B": "b"}, map[string]interface{}{"A": "aa", "B": nil}}
	structSlice := []test{test{A: "a", B: "b"}, test{A: "aa", B: nil}}
	interfaceSlice := []interface{}{test{A: "a", B: "b"}, map[string]interface{}{"A": "aa", "B": nil}}
	result := `A,B
a,b
aa,
`

	var b1, b2, b3, b4 bytes.Buffer
	if err := CSV([]string{"A", "B"}, mapSlice, &b1); err != nil {
		fmt.Println(err)
		t.Error("Export map slice source csv failed")
	}
	if err := CSV([]string{"A", "B"}, structSlice, &b2); err != nil {
		fmt.Println(err)
		t.Error("Export struct slice source csv failed")
	}
	if err := CSV(nil, structSlice, &b3); err != nil {
		fmt.Println(err)
		t.Error("Export struct slice source csv failed")
	}
	if err := CSV([]string{"A", "B"}, interfaceSlice, &b4); err != nil {
		fmt.Println(err)
		t.Error("Export interface slice source csv failed")
	}
	if b1.String() != result {
		t.Error("Export map slice source csv result is not except one")
	}
	if b2.String() != result {
		t.Error("Export struct slice source csv result is not except one")
	}
	if b3.String() != result {
		t.Error("Export struct slice source with nil fieldnames csv result is not except one")
	}
	if b4.String() != result {
		t.Error("Export interface slice source csv result is not except one")
	}
}

func TestCSVWithUTF8BOM(t *testing.T) {
	bom := []byte{0xEF, 0xBB, 0xBF}
	result := `A,B
a,b
`
	var b bytes.Buffer
	if err := CSVWithUTF8BOM([]string{"A", "B"}, []interface{}{map[string]string{"A": "a", "B": "b"}}, &b); err != nil {
		fmt.Println(err)
		t.Error("Export csv with utf8bom failed")
	}
	c := b.Bytes()
	if !reflect.DeepEqual(bom, c[:3]) {
		t.Error("Export csv with utf8bom result not contain bom header")
	}
	if string(c[3:]) != result {
		t.Error("Export csv with utf8bom result is not except one")
	}
}

func TestGetStructFieldNames(t *testing.T) {
	if _, err := getStructFieldNames(map[string]string{"test": "test"}); err == nil {
		t.Error("Except error when get fieldnames from map")
	}
	if fieldnames, err := getStructFieldNames(struct{ A, B string }{}); err != nil {
		t.Error("Failed to get fieldnames")
	} else {
		if !reflect.DeepEqual(fieldnames, []string{"A", "B"}) {
			t.Error("Fieldnames result is not except one")
		}
	}
}

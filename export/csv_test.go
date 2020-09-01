package export

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCSV(t *testing.T) {
	type test struct{ A, B interface{} }
	mapString := []interface{}{map[string]interface{}{"A": "a", "B": "b"}, map[string]interface{}{"A": "aa", "B": nil}}
	Struct := []interface{}{test{A: "a", B: "b"}, test{A: "aa", B: nil}}
	result := `A,B
a,b
aa,
`

	var b1, b2 bytes.Buffer
	if CSV([]string{"A", "B"}, mapString, &b1) != nil {
		t.Error("Export map source csv failed")
	}
	if CSV([]string{"A", "B"}, Struct, &b2) != nil {
		t.Error("Export struct source csv failed")
	}
	c1, _ := ioutil.ReadAll(&b1)
	c2, _ := ioutil.ReadAll(&b2)
	if string(c1) != result {
		t.Error("Export map source csv result is not except one")
	}
	if string(c2) != result {
		t.Error("Export struct source csv result is not except one")
	}
}

func TestCSVWithUTF8BOM(t *testing.T) {
	bom := []byte{0xEF, 0xBB, 0xBF}
	result := `A,B
a,b
`
	var b bytes.Buffer
	if CSVWithUTF8BOM([]string{"A", "B"}, []interface{}{map[string]interface{}{"A": "a", "B": "b"}}, &b) != nil {
		t.Error("Export csv with utf8bom failed")
	}
	c, _ := ioutil.ReadAll(&b)
	if !reflect.DeepEqual(bom, c[:3]) {
		t.Error("Export csv with utf8bom result not contain bom header")
	}
	if string(c[3:]) != result {
		t.Error("Export csv with utf8bom result is not except one")
	}
}

package export

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	var b1, b2, b3 bytes.Buffer
	if err := CSV([]string{"A", "B"}, mapSlice, &b1); err != nil {
		fmt.Println(err)
		t.Error("Export map slice source csv failed")
	}
	if err := CSV([]string{"A", "B"}, structSlice, &b2); err != nil {
		fmt.Println(err)
		t.Error("Export struct slice source csv failed")
	}
	if err := CSV([]string{"A", "B"}, interfaceSlice, &b3); err != nil {
		fmt.Println(err)
		t.Error("Export interface slice source csv failed")
	}
	c1, _ := ioutil.ReadAll(&b1)
	c2, _ := ioutil.ReadAll(&b2)
	c3, _ := ioutil.ReadAll(&b3)
	if string(c1) != result {
		t.Error("Export map slice source csv result is not except one")
	}
	if string(c2) != result {
		t.Error("Export struct slice source csv result is not except one")
	}
	if string(c3) != result {
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
	c, _ := ioutil.ReadAll(&b)
	if !reflect.DeepEqual(bom, c[:3]) {
		t.Error("Export csv with utf8bom result not contain bom header")
	}
	if string(c[3:]) != result {
		t.Error("Export csv with utf8bom result is not except one")
	}
}

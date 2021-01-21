package csv

import (
	"reflect"
	"strings"
	"testing"
)

type result struct{ A, B string }

func TestReader(t *testing.T) {
	csv := `A,B
a,b
1,2
`
	rs, err := ReadAll(strings.NewReader(csv))
	if err != nil {
		t.Error(err)
	}

	fields := rs.Fields()
	if !reflect.DeepEqual([]string{"A", "B"}, fields) {
		t.Errorf("expected %v; got %v", []string{"A", "B"}, fields)
	}

	var results []result
	for rs.Next() {
		var result result
		if err := rs.Scan(&result.A, &result.B); err != nil {
			t.Error(err)
		}
		results = append(results, result)
	}
	if !reflect.DeepEqual([]result{{"a", "b"}, {"1", "2"}}, results) {
		t.Errorf("expected %v; got %v", []result{{"a", "b"}, {"1", "2"}}, results)
	}
}

package workers

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
)

func TestSetMax(t *testing.T) {
	SetMax(5)
	if defaultWorkers.Max != 5 {
		t.Error("SetMax result is not except one")
	}
	SetMax(100)
	if defaultWorkers.Max != 100 {
		t.Error("SetMax result is not except one")
	}
}

func TestRunOnSlice(t *testing.T) {
	type test struct {
		char  string
		times int
	}
	Slice := []test{test{"a", 1}, test{"b", 2}, test{"c", 3}}
	result := make([]string, len(Slice))
	if err := RunOnSlice(Slice, func(i int, item interface{}) {
		result[i] = strings.Repeat(item.(test).char, item.(test).times)
	}); err != nil {
		fmt.Println(err)
		t.Error("RunOnSlice workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("RunOnSlice workers result is not except one")
	}
}

func TestRunOnMap(t *testing.T) {
	Map := map[string]int{"a": 1, "b": 2, "c": 3}
	m := &sync.Mutex{}
	var result []string
	if err := RunOnMap(Map, func(k, v interface{}) {
		m.Lock()
		result = append(result, strings.Repeat(k.(string), v.(int)))
		m.Unlock()
	}); err != nil {
		fmt.Println(err)
		t.Error("RunOnMap workers failed")
	}
	sort.Strings(result)
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		fmt.Println(result)
		t.Error("RunOnMap workers result is not except one")
	}
}

func TestRunOnRange(t *testing.T) {
	end := 3
	items := []string{"a", "b", "c"}
	result := make([]string, end)
	if err := RunOnRange(1, end, func(num int) {
		result[num-1] = strings.Repeat(items[num-1], num)
	}); err != nil {
		fmt.Println(err)
		t.Error("RunOnRange workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("RunOnRange workers result is not except one")
	}
}

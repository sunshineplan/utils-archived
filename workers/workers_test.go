package workers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestSetMax(t *testing.T) {
	SetMax(5)
	if defaultWorkers.Max != 5 {
		t.Error("SetMax result is not except one")
	}
}

func TestRun(t *testing.T) {
	type test struct {
		char  string
		times int
	}
	items := []test{test{"a", 1}, test{"b", 2}, test{"c", 3}}
	result := make([]string, len(items))
	if err := Run(items, func(i int, item interface{}) {
		result[i] = strings.Repeat(item.(test).char, item.(test).times)
	}); err != nil {
		fmt.Println(err)
		t.Error("Run workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("Run workers result is not except one")
	}
}

func TestRunRange(t *testing.T) {
	end := 3
	items := []string{"a", "b", "c"}
	result := make([]string, end)
	if err := RunRange(1, end, func(num int) {
		result[num-1] = strings.Repeat(items[num-1], num)
	}); err != nil {
		fmt.Println(err)
		t.Error("RunRange workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("RunRange workers result is not except one")
	}
}
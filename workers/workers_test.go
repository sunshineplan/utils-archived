package workers

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	type test struct {
		char  string
		times int
	}
	slice := []test{test{"a", 1}, test{"b", 2}, test{"c", 3}}
	result := make([]string, len(slice))
	if err := Slice(slice, func(i int, item interface{}) {
		result[i] = strings.Repeat(item.(test).char, item.(test).times)
	}); err != nil {
		fmt.Println(err)
		t.Error("RunOnSlice workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("RunOnSlice workers result is not except one")
	}
}

func TestMap(t *testing.T) {
	m := &sync.Mutex{}
	var result []string
	if err := Map(map[string]int{"a": 1, "b": 2, "c": 3}, func(k, v interface{}) {
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

func TestRange(t *testing.T) {
	end := 3
	items := []string{"a", "b", "c"}
	result := make([]string, end)
	if err := Range(1, end, func(num int) {
		result[num-1] = strings.Repeat(items[num-1], num)
	}); err != nil {
		fmt.Println(err)
		t.Error("RunOnRange workers failed")
	}
	if !reflect.DeepEqual(result, []string{"a", "bb", "ccc"}) {
		t.Error("RunOnRange workers result is not except one")
	}
}

func TestLimit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	limit := rand.Intn(1000) + 51
	var count1, count2, count3 int
	go Range(1, limit, func(_ int) {
		count1++
		for {
			select {}
		}
	})
	go DefaultRange(1, limit, func(_ int) {
		count2++
		for {
			select {}
		}
	})
	workers := rand.Intn(50) + 1
	go New(workers).Range(1, limit, func(_ int) {
		count3++
		for {
			select {}
		}
	})
	time.Sleep(time.Second)
	if count1 != limit {
		t.Errorf("Unlimited workers goroutine number is not %d: %d", limit, count1)
	}
	if count2 != runtime.NumCPU()*2 {
		t.Errorf("Default workers goroutine number is not %d: %d", runtime.NumCPU()*2, count2)
	}
	if count3 != workers {
		t.Errorf("Workers goroutine number is not %d: %d", workers, count3)
	}
}

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

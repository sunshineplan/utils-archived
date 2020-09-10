package workers

import (
	"fmt"
	"reflect"
	"runtime"
)

var defaultWorkers = Workers{Max: runtime.NumCPU() * 2}

// Workers can run jobs concurrently with limit
type Workers struct {
	Max int
}

// New create new workers with max limit
func New(max int) Workers {
	return Workers{Max: max}
}

// SetMax default workers max limit
func SetMax(max int) {
	defaultWorkers.Max = max
}

// Slice run workers on slice without limit
func Slice(Slice interface{}, Runner func(int, interface{})) error {
	return runSlice(0, Slice, Runner)
}

// Map run workers on slice without limit
func Map(Map interface{}, Runner func(interface{}, interface{})) error {
	return runMap(0, Map, Runner)
}

// Range run workers on range without limit
func Range(Start, End int, Runner func(int)) error {
	return runRange(0, Start, End, Runner)
}

// DefaultSlice use default workers run on slice
func DefaultSlice(Slice interface{}, Runner func(int, interface{})) error {
	return defaultWorkers.Slice(Slice, Runner)
}

// DefaultMap use default workers run on slice
func DefaultMap(Map interface{}, Runner func(interface{}, interface{})) error {
	return defaultWorkers.Map(Map, Runner)
}

// DefaultRange use default workers run on range
func DefaultRange(Start, End int, Runner func(int)) error {
	return defaultWorkers.Range(Start, End, Runner)
}

// Slice run workers on slice
func (w Workers) Slice(Slice interface{}, Runner func(int, interface{})) error {
	return runSlice(w.Max, Slice, Runner)
}

// Map run workers on slice
func (w Workers) Map(Map interface{}, Runner func(interface{}, interface{})) error {
	return runMap(w.Max, Map, Runner)
}

// Range run workers on range
func (w Workers) Range(Start, End int, Runner func(int)) error {
	return runRange(w.Max, Start, End, Runner)
}

func runSlice(Limit int, Slice interface{}, Runner func(int, interface{})) error {
	if reflect.TypeOf(Slice).Kind() != reflect.Slice {
		return fmt.Errorf("First argument must be a slice")
	}
	values := reflect.ValueOf(Slice)
	if Limit <= 0 {
		Limit = values.Len()
	}
	c := make(chan bool, Limit)
	for i := 0; i < values.Len(); i++ {
		c <- true
		go func(index int, value interface{}) {
			defer func() { <-c }()
			Runner(index, value)
		}(i, values.Index(i).Interface())
	}
	for i := 0; i < Limit; i++ {
		c <- true
	}
	return nil
}

func runMap(Limit int, Map interface{}, Runner func(interface{}, interface{})) error {
	if reflect.TypeOf(Map).Kind() != reflect.Map {
		return fmt.Errorf("First argument must be a map")
	}
	value := reflect.ValueOf(Map)
	if Limit <= 0 {
		Limit = len(value.MapKeys())
	}
	iter := value.MapRange()
	c := make(chan bool, Limit)
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		c <- true
		go func(k, v interface{}) {
			defer func() { <-c }()
			Runner(k, v)
		}(k.Interface(), v.Interface())
	}
	for i := 0; i < Limit; i++ {
		c <- true
	}
	return nil
}

func runRange(Limit, Start, End int, Runner func(int)) error {
	if Limit <= 0 {
		Limit = End - Start + 1
	}
	c := make(chan bool, Limit)
	for i := Start; i <= End; i++ {
		c <- true
		go func(num int) {
			defer func() { <-c }()
			Runner(num)
		}(i)
	}
	for i := 0; i < Limit; i++ {
		c <- true
	}
	return nil
}

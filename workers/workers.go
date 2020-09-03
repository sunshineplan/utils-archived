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

// RunOnSlice default workers on slice
func RunOnSlice(Slice interface{}, Runner func(int, interface{})) error {
	return defaultWorkers.RunOnSlice(Slice, Runner)
}

// RunOnMap default workers on slice
func RunOnMap(Map interface{}, Runner func(interface{}, interface{})) error {
	return defaultWorkers.RunOnMap(Map, Runner)
}

// RunOnRange run default workers on range
func RunOnRange(Start, End int, Runner func(int)) error {
	return defaultWorkers.RunOnRange(Start, End, Runner)
}

// RunOnSlice workers on slice
func (w Workers) RunOnSlice(Slice interface{}, Runner func(int, interface{})) error {
	if reflect.TypeOf(Slice).Kind() != reflect.Slice {
		return fmt.Errorf("First argument must be a slice")
	}
	values := reflect.ValueOf(Slice)
	c := make(chan bool, w.Max)
	for i := 0; i < values.Len(); i++ {
		c <- true
		go func(index int, value interface{}) {
			defer func() { <-c }()
			Runner(index, value)
		}(i, values.Index(i).Interface())
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

// RunOnMap workers on slice
func (w Workers) RunOnMap(Map interface{}, Runner func(interface{}, interface{})) error {
	if reflect.TypeOf(Map).Kind() != reflect.Map {
		return fmt.Errorf("First argument must be a map")
	}
	iter := reflect.ValueOf(Map).MapRange()
	c := make(chan bool, w.Max)
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		c <- true
		go func(k, v interface{}) {
			defer func() { <-c }()
			Runner(k, v)
		}(k.Interface(), v.Interface())
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

// RunOnRange run workers on range
func (w Workers) RunOnRange(Start, End int, Runner func(int)) error {
	c := make(chan bool, w.Max)
	for i := Start; i <= End; i++ {
		c <- true
		go func(num int) {
			defer func() { <-c }()
			Runner(num)
		}(i)
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

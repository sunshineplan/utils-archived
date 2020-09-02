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

// Run default workers on slice
func Run(slice interface{}, runner func(int, interface{})) error {
	return defaultWorkers.Run(slice, runner)
}

// RunRange run default workers on range
func RunRange(start, end int, runner func(int)) error {
	return defaultWorkers.RunRange(start, end, runner)
}

// Run workers on slice
func (w Workers) Run(slice interface{}, runner func(int, interface{})) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return fmt.Errorf("First argument must be a slice")
	}
	values := reflect.ValueOf(slice)
	c := make(chan bool, w.Max)
	for i := 0; i < values.Len(); i++ {
		c <- true
		go func(index int, value interface{}) {
			defer func() { <-c }()
			runner(index, value)
		}(i, values.Index(i).Interface())
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

// RunRange run workers on range
func (w Workers) RunRange(start, end int, runner func(int)) error {
	c := make(chan bool, w.Max)
	for i := start; i <= end; i++ {
		c <- true
		go func(num int) {
			defer func() { <-c }()
			runner(num)
		}(i)
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

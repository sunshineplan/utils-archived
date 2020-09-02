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
func New(max int) *Workers {
	return &Workers{Max: max}
}

// SetMax default workers max limit
func SetMax(max int) {
	defaultWorkers = Workers{Max: max}
}

// Run default workers on slice
func Run(slice interface{}, runner func(chan bool, int, interface{})) error {
	return defaultWorkers.Run(slice, runner)
}

// RunRange run default workers on range
func RunRange(start, end int, runner func(chan bool, int)) error {
	return defaultWorkers.RunRange(start, end, runner)
}

// Run workers on slice
func (w *Workers) Run(slice interface{}, runner func(chan bool, int, interface{})) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return fmt.Errorf("Input item must be a slice")
	}
	items := reflect.ValueOf(slice)
	c := make(chan bool, w.Max)
	for i := 0; i < items.Len(); i++ {
		c <- true
		go runner(c, i, items.Index(i).Interface())
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

// RunRange run workers on range
func (w *Workers) RunRange(start, end int, runner func(chan bool, int)) error {
	c := make(chan bool, w.Max)
	for i := start; i <= end; i++ {
		c <- true
		go runner(c, i)
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
	return nil
}

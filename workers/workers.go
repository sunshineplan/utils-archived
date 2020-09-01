package workers

import (
	"fmt"
	"reflect"
)

// Workers can run jobs concurrently with limit
type Workers struct {
	Max int
}

// New create new Workers with max limit
func New(max int) *Workers {
	return &Workers{Max: max}
}

// Run workers
func (w *Workers) Run(itemSlice interface{}, runner func(chan bool, int, interface{})) error {
	if reflect.TypeOf(itemSlice).Kind() != reflect.Slice {
		return fmt.Errorf("Input item must be a slice")
	}
	items := reflect.ValueOf(itemSlice)
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

// RunRange run workers according range
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

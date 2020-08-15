package workers

// Workers can run jobs concurrently with limit
type Workers struct {
	Max int
}

// New create new Workers with max limit
func New(max int) *Workers {
	return &Workers{Max: max}
}

// Run workers
func (w *Workers) Run(items []interface{}, runner func(chan bool, int, interface{})) {
	c := make(chan bool, w.Max)
	for i, item := range items {
		c <- true
		go runner(c, i, item)
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
}

// RunRange run workers according range
func (w *Workers) RunRange(start, end int, runner func(chan bool, int)) {
	c := make(chan bool, w.Max)
	for i := start; i <= end; i++ {
		c <- true
		go runner(c, i)
	}
	for i := 0; i < w.Max; i++ {
		c <- true
	}
}

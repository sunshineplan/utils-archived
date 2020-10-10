package retry

import "time"

// Do keeps retrying the function until no error is returned.
func Do(fn func() error, attempts, delay uint) (err error) {
	var i uint
	for ; i < attempts; i++ {
		if err = fn(); err == nil {
			return
		}
		if i < attempts-1 {
			time.Sleep(time.Second * time.Duration(delay))
		}
	}
	return
}

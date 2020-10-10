package retry

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	if err := Do(func() error {
		return nil
	}, 3, 1); err != nil {
		t.Error("test except non error")
	}
	if err := Do(func() error {
		return errors.New("error")
	}, 3, 1); err == nil {
		t.Error("test except error")
	}
}

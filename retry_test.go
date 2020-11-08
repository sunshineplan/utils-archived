package utils

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	if err := Retry(func() error {
		return nil
	}, 3, 1); err != nil {
		t.Error(err)
	}
	if err := Retry(func() error {
		return errors.New("error")
	}, 3, 1); err == nil {
		t.Error("gave nil error; want error")
	}
}

package agent

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	agent, err := Get()
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(agent, "Chrome") {
		t.Error("Get chrome agent failed")
	}
}

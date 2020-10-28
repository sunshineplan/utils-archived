package ocr

import (
	"os"
	"testing"
)

func TestOCR(t *testing.T) {
	f, err := os.Open("ocr.space.logo.png")
	if err != nil {
		t.Error(err)
	}
	r, err := Read(f)
	if err != nil {
		t.Error(err)
	}
	if r != "OCR .Space\r\n" {
		t.Errorf("expected %q; got %q", "OCR .Space\r\n", r)
	}
}

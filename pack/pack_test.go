package pack

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var files = []File{
	{Name: "testdata/1.txt", Body: []byte("1")},
	{Name: "testdata/2.txt", Body: []byte("2")},
}

func TestPack(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	if err := FromBytes(&buf1, ZIP, files...); err != nil {
		t.Error(err)
	}
	if err := FromFiles(&buf2, ZIP, "testdata/1.txt", "testdata/2.txt"); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Error("expected equal zip archive; got not equal")
	}
	if err := FromBytes(ioutil.Discard, TAR, files...); err != nil {
		t.Error(err)
	}
	if err := FromFiles(ioutil.Discard, TAR, "testdata/1.txt", "testdata/2.txt"); err != nil {
		t.Error(err)
	}
}

package zip

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"path/filepath"
)

// File struct contains bytes body and the provided name field.
type File struct {
	Name string
	Body []byte
}

// FromBytes creates a zip archive from bytes and using the provided name.
func FromBytes(w io.Writer, files ...File) error {
	zw := zip.NewWriter(w)

	for _, file := range files {
		f, err := zw.Create(file.Name)
		if err != nil {
			return err
		}
		if _, err := f.Write(file.Body); err != nil {
			return err
		}
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return nil
}

// FromFile creates a zip archive from files and using the base filename.
func FromFile(w io.Writer, files ...string) error {
	zw := zip.NewWriter(w)

	for _, file := range files {
		f, err := zw.Create(filepath.Base(file))
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := f.Write(b); err != nil {
			return err
		}
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return nil
}

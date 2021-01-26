package pack

import (
	"archive/zip"
	"io"
	"io/ioutil"
)

func zipBytes(w io.Writer, files ...File) error {
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

func zipFiles(w io.Writer, files ...string) error {
	zw := zip.NewWriter(w)

	for _, file := range files {
		f, err := zw.Create(file)
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

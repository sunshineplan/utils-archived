package unpack

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Zip unpacks zip archive
func Zip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			return err
		}

		if err := outFile.Close(); err != nil {
			return err
		}
		if err := rc.Close(); err != nil {
			return err
		}

	}
	return nil
}

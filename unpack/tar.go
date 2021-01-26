package unpack

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Tar unpacks tar.gz archive
func Tar(r io.Reader, dest string) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(filepath.Join(dest, header.Name)); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(filepath.Join(dest, header.Name), 0755); err != nil {
						return err
					}
				} else {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.Create(filepath.Join(dest, header.Name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}

		default:
			log.Printf(
				"ExtractTarGz: uknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}
	}

	return nil
}

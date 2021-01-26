package pack

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func tarBytes(w io.Writer, files ...File) error {
	gw := gzip.NewWriter(w)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		header := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := tw.Write(file.Body); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	return nil
}

func tarFiles(w io.Writer, files ...string) error {
	gw := gzip.NewWriter(w)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = file

		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	return nil
}

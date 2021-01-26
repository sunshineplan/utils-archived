package pack

import (
	"errors"
	"io"
)

// File struct contains bytes body and the provided name field.
type File struct {
	Name string
	Body []byte
}

// Format represents the archive format.
type Format int

const (
	// ZIP format
	ZIP Format = iota
	// TAR format
	TAR
)

// FromBytes creates an archive from bytes and using the provided name.
func FromBytes(w io.Writer, format Format, files ...File) error {
	switch format {
	case ZIP:
		return zipBytes(w, files...)
	case TAR:
		return tarBytes(w, files...)
	default:
		return errors.New("Unknow format")
	}
}

// FromFiles creates an archive from files and using the base filename.
func FromFiles(w io.Writer, format Format, files ...string) error {
	switch format {
	case ZIP:
		return zipFiles(w, files...)
	case TAR:
		return tarFiles(w, files...)
	default:
		return errors.New("unknow format")
	}
}

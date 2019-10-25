package v2

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	of "github.com/cisco-cx/of/pkg/v2"
)

// Implements Concatenate.
type Files struct {
	Path string // Dir containing files to concat.
	Ext  string // Extention of file. If not all files in `Path` will be concatenated
}

// Concatenate all files into
func (f *Files) Concat() (io.Reader, error) {

	//Check if given path is a directory.
	fi, err := os.Stat(f.Path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() == false {
		return nil, of.ErrPathIsNotDir
	}

	// Create path to be Globbed
	path := ""
	if f.Ext != "" {
		path = filepath.Join(f.Path, fmt.Sprintf("*.%s", f.Ext))
	} else {
		path = filepath.Join(f.Path, "*")
	}

	// Get files in path.
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	// Add each file to the io.Reader buffer.
	data := bytes.NewBuffer(nil)
	for _, f := range files {
		handle, err := os.Open(f)
		if err != nil {
			return data, err
		}
		defer handle.Close()

		_, err = data.ReadFrom(handle)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

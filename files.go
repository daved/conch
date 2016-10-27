package main

import (
	"io/ioutil"
	"os"
	"strings"
)

// fileOutput holds a file path (w/ dir), processed data, and error (if any).
type fileOutput struct {
	path string
	data string
	err  error
}

// fileInfoGroup holds a slice of os.FileInfo along with the directory that the
// info came from.
type fileInfoGroup struct {
	dir string
	fsi []os.FileInfo
}

// newFileInfoGroup returns a pointer to a fileInfoGroup populated by all
// gzipped files within the provided directory with a depth of 1.
func newFileInfoGroup(dir string) (*fileInfoGroup, error) {
	fsi, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for k := len(fsi) - 1; k >= 0; k-- {
		// remove directories, and files without the correct extension
		if fsi[k].IsDir() || !strings.HasSuffix(fsi[k].Name(), ".gz") {
			fsi = append(fsi[:k], fsi[k+1:]...)
		}
	}

	return &fileInfoGroup{dir: dir, fsi: fsi}, nil
}

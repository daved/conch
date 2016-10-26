package main

import (
	"io/ioutil"
	"os"
	"strings"
)

// file holds a full file path, processed data, and error (if any).
type file struct {
	path string
	data string
	err  error
}

// filesInfo holds a slice of os.FileInfo along with the directory the
// contents came from.
type filesInfo struct {
	dir string
	fsi []os.FileInfo
}

// gatherFilesInfo is a convenience function which grabs all gzipped files
// within the provided directory with a depth of 1.
func gatherFilesInfo(dir string) (*filesInfo, error) {
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

	return &filesInfo{dir: dir, fsi: fsi}, nil
}

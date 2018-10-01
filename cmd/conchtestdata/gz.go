package main

import (
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func digitLength(i int) int {
	for n := 1; n < 100000; n++ {
		if int64(float64(i)/math.Pow10(n)) == 0 {
			return n
		}
	}

	return 0
}

func createGZFile(dir string, dLen, i int) error {
	wrapError := func(err error) error {
		return fmt.Errorf("cannot create gz file: %s", err)
	}

	padNum := fmt.Sprintf("%0"+strconv.Itoa(dLen)+"d", i)
	filename := fmt.Sprintf("file%s.gz", padNum)
	filedata := []byte(fmt.Sprintf("This is file #%s with junk content.\n", padNum))

	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return wrapError(err)
	}
	defer f.Close() //nolint

	gzw := gzip.NewWriter(f)
	if _, err = gzw.Write(filedata); err != nil {
		return wrapError(err)
	}
	defer gzw.Close() //nolint

	return nil
}

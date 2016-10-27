package main

import (
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func digitLen(i int) int {
	for n := 1; n < 100000; n++ {
		if int64(float64(i)/math.Pow10(n)) == 0 {
			return n
		}
	}

	return 0
}

func createGZFile(dir string, max, i int) error {
	dLen := digitLen(max)
	padNum := fmt.Sprintf("%0"+strconv.Itoa(dLen)+"d", i)

	filename := fmt.Sprintf("file%s.gz", padNum)
	filedata := []byte(fmt.Sprintf("This is file #%s.\n", padNum))

	f, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		_ = e
	}()

	gzw := gzip.NewWriter(f)
	if _, err = gzw.Write(filedata); err != nil {
		return err
	}
	defer func() {
		e := gzw.Close()
		_ = e
	}()

	return err
}

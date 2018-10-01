package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err) //nolint
		os.Exit(1)
	}
}

func run() error {
	var (
		qty = 1024
		dir = "./testdata"
	)
	flag.IntVar(&qty, "qty", qty, `quantity of test files`)
	flag.StringVar(&dir, "dir", dir, `location of test files`)
	flag.Parse()

	if qty < 1 {
		qty = 1
	}
	dir = filepath.Clean(dir)

	if err := prepDir(dir); err != nil {
		return err
	}

	dLen := digitLength(qty)
	for i := 0; i < qty; i++ {
		if err := createGZFile(dir, dLen, i); err != nil {
			return err
		}
	}

	return nil
}

func prepDir(dir string) error {
	wrapError := func(err error) error {
		return fmt.Errorf("cannot prepare directory: %s", err)
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		if err := os.RemoveAll(dir); err != nil {
			return wrapError(err)
		}
	}

	if err := os.Mkdir(dir, 0775); err != nil { //nolint
		return wrapError(err)
	}

	return nil
}

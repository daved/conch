package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type options struct {
	qty int
	dir string
}

func newOptionsDefault() *options {
	return &options{
		qty: 1024,
		dir: "./testfiles",
	}
}

func main() {
	opts := newOptionsDefault()
	{
		flag.IntVar(&opts.qty, "qty", opts.qty,
			`quantity of test files`)
		flag.StringVar(&opts.dir, "dir", opts.dir,
			`location of test files`)
	}
	flag.Parse()

	if opts.qty < 1 {
		opts.qty = 1
	}

	opts.dir = filepath.Clean(opts.dir)

	if _, err := os.Stat(opts.dir); !os.IsNotExist(err) {
		if err := os.RemoveAll(opts.dir); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	if err := os.Mkdir(opts.dir, 0700); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for i := 0; i < opts.qty; i++ {
		if err := createGZFile(opts.dir, opts.qty, i); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
}

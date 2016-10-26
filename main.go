package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// slow enables a slowing of the digest function and may help users to
	// better understand the implementation of concurrency
	slow = false

	// width controls the amount of goroutines running the digest function
	width = 8
)

func main() {
	// define and parse flags
	{
		flag.BoolVar(&slow, "slow", slow,
			`Slow processing to clarify behavior.`)
		flag.IntVar(&width, "width", width,
			`Set concurrency width.`)
	}
	flag.Parse()

	// get populated filesInfo type
	fsi, err := gatherFilesInfo("./testdata")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// setup done channel and clean up
	done := make(chan struct{})
	defer close(done)

	// get files and error channels
	fs, errc := funnel(done, fsi)

	// print file contents or error
	for f := range fs {
		select {
		case err := <-errc:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		default:
			fmt.Println(f.path, f.data, f.err)
		}
	}
}

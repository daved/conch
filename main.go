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
	flag.BoolVar(&slow, "slow", slow, `slow processing to clarify behavior`)
	flag.IntVar(&width, "width", width, `set concurrency width`)
	flag.Parse()

	// get fileInfoGroup
	fig, err := newFileInfoGroup("./testfiles")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// get new conch and setup cleanup
	c := newConch(fig)
	defer close(c.doneChan())

	// get fileOutput and error channels
	fos, errc := c.run()

	// print file contents or error
	for fo := range fos {
		select {
		case err := <-errc:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		default:
			fmt.Println(fo.path, fo.data, fo.err)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codemodus/sigmon"
)

var (
	// slow enables a slowing of the digest function and may help users to
	// better understand the implementation of concurrency
	slow = false

	// width controls the amount of goroutines running the digest function
	width = 8
)

func main() {
	sm := sigmon.New(nil)
	sm.Run()

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

	sm.Set(func(s *sigmon.SignalMonitor) {
		close(c.done())
		sm.Stop()
	})

	// get fileOutput and error channels
	fos, errs := c.run()

	// print file contents
	for fo := range fos {
		fmt.Println(fo.path, fo.data, fo.err)
	}

	// print error, if any
	select {
	case err := <-errs:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	default:
	}

	sm.Stop()
}

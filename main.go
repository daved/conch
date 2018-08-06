package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codemodus/sigmon"
)

func main() {
	sm := sigmon.New(nil)
	sm.Start()

	var (
		slow  = false
		width = 8
	)
	flag.BoolVar(&slow, "slow", slow, `slow processing to clarify behavior`)
	flag.IntVar(&width, "width", width, `set concurrency width`)
	flag.Parse()

	paths, err := gzipFilePaths("./testfiles")
	if err != nil {
		fmt.Printf("cannot get paths: %s\n", err)
		os.Exit(1)
	}
	c := newConch()

	sm.Set(func(s *sigmon.State) {
		close(c.done()) // trip on any system signal
	})

	fis, runErr := c.run(slow, width, paths)
	for fi := range fis {
		fmt.Println(fi.path, fi.data, fi.err)
	}

	if err = runErr(); err != nil {
		fmt.Printf("run error: %s\n", err)
	}

	sm.Stop()
}

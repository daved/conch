package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codemodus/sigmon"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err) //nolint
		os.Exit(1)
	}
}

func run() error {
	var (
		slow  = false
		width = 8
	)
	flag.BoolVar(&slow, "slow", slow, `slow processing to clarify behavior`)
	flag.IntVar(&width, "width", width, `set concurrency width`)
	flag.Parse()

	paths, err := gzipFilePaths("./testfiles")
	if err != nil {
		return err
	}
	done := make(chan struct{})
	defer safeClose(done)

	sm := sigmon.New(func(s *sigmon.State) {
		safeClose(done) // trip on any system signal
	})
	sm.Start()
	defer sm.Stop()

	fis, fisErr := fileInfos(done, slow, width, paths)
	for fi := range fis {
		fmt.Println(fi.path, fi.data, fi.err)
	}

	return fisErr()
}

func safeClose(c chan struct{}) {
	select {
	case <-c:
	default:
		close(c)
	}
}

package main

import (
	"errors"
	"path/filepath"
	"sync"
)

type conch struct {
	fig  *fileInfoGroup
	done chan struct{}
	out  chan fileOutput
	err  chan error
}

// newConch returns a pointer to a new conch with channels setup.
func newConch(fig *fileInfoGroup) *conch {
	return &conch{
		fig:  fig,
		done: make(chan struct{}),
		out:  make(chan fileOutput),
		err:  make(chan error, 1),
	}
}

// doneChan returns the conch done channel.
func (c *conch) doneChan() chan struct{} {
	return c.done
}

// feedPaths is a generator that sends paths out to digest goroutines.
func (c *conch) feedPaths(paths chan string) {
	for _, v := range c.fig.fsi {
		select {
		case paths <- filepath.Join(c.fig.dir, v.Name()):
		case <-c.done:
			c.err <- errors.New("canceled")
		}
	}

	close(paths)
}

// fanout sets-up digest goroutines according to width.
func (c *conch) fanout(paths chan string) {
	var wg sync.WaitGroup
	wg.Add(width)

	// setup digesters by width
	for i := 0; i < width; i++ {
		go func() {
			digest(c.done, paths, c.out)
			wg.Done()
		}()
	}

	wg.Wait()
	close(c.out)
}

// run wires up the paths chan, calls the path generator, calls the fanout, and
// returns the out and err channels.
func (c *conch) run() (<-chan fileOutput, <-chan error) {
	paths := make(chan string)

	go c.feedPaths(paths)
	go c.fanout(paths)

	return c.out, c.err
}

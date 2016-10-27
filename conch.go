package main

import (
	"errors"
	"path/filepath"
	"sync"
	"time"
)

// conch holds a file info group, and done, out, and err channels.
type conch struct {
	fig  *fileInfoGroup
	done chan struct{}
	out  chan *fileOutput
	err  chan error
}

// newConch returns a pointer to a new conch with channels setup.
func newConch(fig *fileInfoGroup) *conch {
	return &conch{
		fig:  fig,
		done: make(chan struct{}),
		out:  make(chan *fileOutput),
		err:  make(chan error, 1),
	}
}

// doneChan returns the conch done channel.
func (c *conch) doneChan() chan struct{} {
	return c.done
}

// feedPaths is a generator that sends paths out for processing. If canceled,
// an error is returned.
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

// digest processes the file located at the currently provided path, and
// sends out a new result.
func (c *conch) digest(paths <-chan string) {
	for p := range paths {
		r := newFileOutput(p)

		if slow {
			select {
			case <-time.After(time.Second):
			case <-c.done:
				return
			}
		}

		select {
		case c.out <- r:
		case <-c.done:
			return
		}
	}
}

// fanout sets-up digest goroutines according to width. Each digest goroutine
// waits for data from the paths channel and the entire function collapses when
// completed.
func (c *conch) fanout(paths chan string) {
	var wg sync.WaitGroup
	wg.Add(width)

	// setup digesters by width
	for i := 0; i < width; i++ {
		go func() {
			c.digest(paths)
			wg.Done()
		}()
	}

	wg.Wait()
	close(c.out)
}

// run wires up the paths chan, calls the path generator, calls the fanout, and
// returns the out and err channels.
func (c *conch) run() (<-chan *fileOutput, <-chan error) {
	paths := make(chan string)

	go c.feedPaths(paths)
	go c.fanout(paths)

	return c.out, c.err
}

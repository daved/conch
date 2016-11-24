package main

import (
	"errors"
	"path/filepath"
	"sync"
	"time"
)

// conch holds a file info group, and done channel.
type conch struct {
	fig *fileInfoGroup
	d   chan struct{}
}

// newConch returns a pointer to a new conch with channel setup.
func newConch(fig *fileInfoGroup) *conch {
	return &conch{
		fig: fig,
		d:   make(chan struct{}),
	}
}

// done returns the done channel.
func (c *conch) done() chan struct{} {
	return c.d
}

// produce is a generator that produces paths for processing. If canceled, an
// error is produced.
func (c *conch) produce() (<-chan string, <-chan error) {
	paths := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(paths)

		for _, v := range c.fig.fsi {
			select {
			case paths <- filepath.Join(c.fig.dir, v.Name()):
			case <-c.d:
				errs <- errors.New("canceled")
				return
			}
		}
	}()

	return paths, errs
}

// digest processes the file located at the currently provided path, and sends
// out a new result.
func (c *conch) digest(paths <-chan string, fos chan<- *fileOutput) {
	for p := range paths {
		fo := newFileOutput(p)

		if slow {
			select {
			case <-time.After(time.Second):
			case <-c.d:
				return
			}
		}

		select {
		case fos <- fo:
		case <-c.d:
			return
		}
	}
}

// consume sets-up digest goroutines according to width. Each digest goroutine
// waits for data from the paths channel and the entire function collapses when
// completed.
func (c *conch) consume(paths <-chan string) <-chan *fileOutput {
	fos := make(chan *fileOutput)

	go func() {
		var wg sync.WaitGroup
		wg.Add(width)

		for i := 0; i < width; i++ {
			go func() {
				c.digest(paths, fos)
				wg.Done()
			}()
		}

		wg.Wait()
		close(fos)
	}()

	return fos
}

// run calls the path generator, calls the consume func using the returned
// paths channel, and returns the outs and errs channels.
func (c *conch) run() (<-chan *fileOutput, <-chan error) {
	paths, errs := c.produce()
	fos := c.consume(paths)

	return fos, errs
}

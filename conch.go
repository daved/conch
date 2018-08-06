package main

import (
	"errors"
	"sync"
	"time"
)

type conch struct {
	d chan struct{}
}

func newConch() *conch {
	return &conch{
		d: make(chan struct{}),
	}
}

func (c *conch) done() chan struct{} {
	return c.d
}

func (c *conch) produce(paths []string) (<-chan string, <-chan error) {
	psc := make(chan string)
	ec := make(chan error, 1)

	go func() {
		defer close(psc)

		for _, p := range paths {
			select {
			case psc <- p:
			case <-c.d:
				ec <- errors.New("canceled")
				return
			}
		}
	}()

	return psc, ec
}

func (c *conch) digest(slow bool, fisc chan<- *fileInfo, psc <-chan string) {
	for p := range psc {
		if slow {
			select {
			case <-time.After(time.Second):
			case <-c.d:
				return
			}
		}

		fi := newFileInfo(p)
		select {
		case fisc <- fi:
		case <-c.d:
			return
		}
	}
}

func (c *conch) consume(slow bool, width int, psc <-chan string) <-chan *fileInfo {
	fisc := make(chan *fileInfo)

	go func() {
		var wg sync.WaitGroup
		wg.Add(width)

		for i := 0; i < width; i++ {
			go func() {
				c.digest(slow, fisc, psc)
				wg.Done()
			}()
		}

		wg.Wait()
		close(fisc)
	}()

	return fisc
}

func (c *conch) run(slow bool, width int, paths []string) (<-chan *fileInfo, func() error) {
	psc, ec := c.produce(paths)
	fisc := c.consume(slow, width, psc)

	errFn := func() error {
		select {
		case err := <-ec:
			return err
		default:
			return nil
		}
	}

	return fisc, errFn
}

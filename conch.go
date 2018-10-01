package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func produce(done <-chan struct{}, paths []string) (<-chan string, <-chan error) {
	psc := make(chan string)
	esc := make(chan error, 1)

	go func() {
		defer close(psc)
		defer close(esc)

		for _, p := range paths {
			select {
			case psc <- p:
			case <-done:
				esc <- errors.New("canceled")
				return
			}
		}
	}()

	return psc, esc
}

func digest(done <-chan struct{}, slow bool, fisc chan<- *fileInfo, psc <-chan string) {
	for p := range psc {
		if slow {
			select {
			case <-time.After(time.Second):
			case <-done:
				return
			}
		}

		select {
		case fisc <- newFileInfo(p):
		case <-done:
			return
		}
	}
}

func consume(done <-chan struct{}, slow bool, width int, psc <-chan string) <-chan *fileInfo {
	fisc := make(chan *fileInfo)

	go func() {
		defer close(fisc)

		var wg sync.WaitGroup
		wg.Add(width)

		for i := 0; i < width; i++ {
			go func() {
				digest(done, slow, fisc, psc)
				wg.Done()
			}()
		}

		wg.Wait()
	}()

	return fisc
}

func fileInfos(done <-chan struct{}, slow bool, width int, paths []string) (<-chan *fileInfo, func() error) {
	psc, esc := produce(done, paths)
	fisc := consume(done, slow, width, psc)

	var last error
	errFn := func() error {
		if err := <-esc; err != nil {
			last = fmt.Errorf("cannot handle fileInfos: %s", err)
			return last
		}

		return last
	}

	return fisc, errFn
}

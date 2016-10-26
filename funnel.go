package main

import (
	"errors"
	"path"
	"sync"
)

// funnel receives a filesInfo type, spawns goroutines (determined by the
// const "width"), and closes
func funnel(done <-chan struct{}, fsi *filesInfo) (<-chan file, <-chan error) {
	c := make(chan file)
	errc := make(chan error, 1)

	var wg sync.WaitGroup
	paths := make(chan string)

	go func() {
		// anon go func sends paths down the correct channel.
		go func() {
			for _, v := range fsi.fsi {
				select {
				case paths <- path.Join(fsi.dir, v.Name()):
				case <-done:
					errc <- errors.New("canceled")
				}
			}

			close(paths)
		}()

		// setup digesters by width.
		wg.Add(width)
		for i := 0; i < width; i++ {
			go func() {
				digest(done, paths, c)
				wg.Done()
			}()
		}

		// wait and close result channel after paths have been processed.
		go func() {
			wg.Wait()
			close(c)
		}()
	}()

	return c, errc
}

package main

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"time"
)

// digest processes the file located at the currently provided path, and
// sends out a new result. It could be used, instead, to communicate with
// relevant micorservices.
func digest(done <-chan struct{}, paths <-chan string, c chan<- file) {
	for p := range paths {
		r := file{path: p}

		func() {
			f, err := os.Open(p)
			if err != nil {
				r.err = err
				return
			}
			defer func() {
				_ = f.Close()
			}()

			gzr, err := gzip.NewReader(f)
			if err != nil {
				r.err = err
				return
			}
			defer func() {
				_ = gzr.Close()
			}()

			data, err := ioutil.ReadAll(gzr)
			if err != nil {
				r.err = err
				return
			}

			r.data = string(data)
		}()

		if slow {
			time.Sleep(time.Second)
		}

		select {
		case c <- r:
		case <-done:
			return
		}
	}
}

package main

import (
	"flag"
	"testing"
)

func TestTrace(t *testing.T) {
	flag.VisitAll(func(fl *flag.Flag) {
		if fl.Name == "test.trace" && fl.Value.String() == "" {
			t.SkipNow()
		}
	})

	main()
}

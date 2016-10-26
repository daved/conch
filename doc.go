/*

conch0 collects file info from the "testdata" directory and concurrently
processes the files by decompressing the contents and then printing the
data or related error.

The "width" of concurrency is set by the constant "width". Parallelism is
scheduled properly regardless of CPUs available, and the processing will
be serial if only one CPU is available. Width, in this case, helps control
the maximum available goroutines to limit the usage of RAM (see heap
profile results).

Usage:
	* This is not properly setup to be built. Use "go run main.go".

Available flags:
	-slow
		Slow processing to clarify behavior.
	-profmem={filename}
		Run memory profile and write to named file.

*/
package main

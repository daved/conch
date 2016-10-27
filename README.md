# conch0

    go get -u github.com/daved/conch0/...

conch0 collects file info from the "testfiles" directory and concurrently
processes the files by decompressing the contents and then printing the data or
related error.

The "width" of concurrency is set by the constant "width". Parallelism is
scheduled properly regardless of CPUs available, and the processing will be
serial if only one CPU is available. Width, in this case, helps control the
maximum available goroutines to limit the usage of RAM (see heap profile
results).

    Available flags:
        -slow
    		slow processing to clarify behavior
    	-width int
    		set concurrency width (default 8)

For convenience, a sub-command has been provided (conchtestdata) which will
generate the required files for processing.

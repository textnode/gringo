gringo
============

A high-performance minimalist queue implemented using a stripped-down lock-free ringbuffer, written in Go (golang.org)

When operating with 2 or more goroutines, GOMAXPROCS >= number of goroutines and sufficient CPU cores to service the goroutines in parallel, 
this gives approximately 6 times the throughput of an equivalent pipeline built using channels.

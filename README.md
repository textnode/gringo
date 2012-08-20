gringo
============

A high-performance minimalist queue implemented using a stripped-down lock-free ringbuffer, written in Go (golang.org)

When operating with 2 or more goroutines (and GOMAXPROCS >= number of goroutines) this gives approximately 6 times the performance of an equivalent pipeline built using channels.

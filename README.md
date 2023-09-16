# GO Server for simpleinsurance challenge

## Goal

Using only the standard library, create a Go HTTP server that on each request responds with a
counter of the total number of requests that it has received during the previous 60 seconds
(moving window). The server should continue to return the correct numbers after restarting it, by
persisting data to a file.

## TODO

1. Actually implement the moving window
2. Tests
3. Clean up code
4. Add Makefile to run code + tests
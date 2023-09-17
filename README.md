# GO Server for simpleinsurance challenge

## Goal

Using only the standard library, create a Go HTTP server that on each request responds with a
counter of the total number of requests that it has received during the previous 60 seconds
(moving window). The server should continue to return the correct numbers after restarting it, by
persisting data to a file.


## Application

Since the requirements were pretty strict, I created a simple HTTP server with one endpoint `/counter` that has no arguments and returns a Dictionary with the key `counter` signifying the number of requests that were made inside the window.

I did not count the request that gives this information as part of the counter, because it arrives at the same timeas we are collecting the data.
In any case, I think both cases are valid but choose this one.

## Running the Application

To run the program either call 

1. Run the main file:  `go run main.go` 

or 

1. Build the app:  `go build .`
2. Run the app:  `./counter`

And then use `curl` to make the requests

`curl http://localhost:8080/counter`


## Testing

To run the full test suite please run: 

`go test -cover  -v ./...`

The 

## TODO

2. Tests
3. Clean up code
4. Add Makefile to run code + tests
5. Split tests that are covering valid/invalid in the same test case
6. Remove test.log after test is done
7. Tests to main and handlers
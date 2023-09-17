package main

import (
	"fmt"
	"log"
	"net/http"
	"simpleinsurance/handlers"
	"simpleinsurance/services"
	"sync"
)

const address string = "localhost:8080"
const windowSize int = 60

// Make this variable Upper case to access it on the handlers file when we have it
var RequestsCounter int
var CounterID int
var lock sync.Mutex

var logFilePath string = "requests.log"

func startup() error {
	exists, err := services.InitFile(logFilePath)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if exists {
		RequestsCounter, _, err := services.ParseLogFile(logFilePath, windowSize)

		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Println(RequestsCounter)
	}

	return nil
}

func Counter(w http.ResponseWriter, r *http.Request) {
	response, err := handlers.UpdateCounter(w, logFilePath, windowSize)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(*response)
}

func main() {
	fmt.Println("Starting Server")

	err := startup()
	if err != nil {
		log.Fatal("Server has encountered a problem")
		panic(err)
	}

	handler := http.HandlerFunc(Counter)
	http.Handle("/counter", handler)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server has encountered a problem")
		panic(err)
	}
}

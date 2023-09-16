package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simpleinsurance/services"
	"simpleinsurance/values"
	"sync"
	"time"
)

const address string = "localhost:8080"
const windowSize int = 60

// Make this variable Upper case to access it on the handlers file when we have it
var RequestsCounter int
var CounterID int
var lock sync.Mutex

var logFilePath string = "logs.txt"

func UpdateCounter(w http.ResponseWriter) (*[]byte, error) {
	lock.Lock()

	timestamp := time.Now().Local().Format(time.RFC3339)

	RequestsCounter, CounterID, err := services.ReadLogLines(logFilePath, windowSize)

	var logLine values.LogLine
	logLine.ID = CounterID + 1
	logLine.Timestamp = timestamp

	fmt.Println(logLine)

	services.WriteToFile(logFilePath, &logLine)

	// Create a Response value to output to the user
	response := values.Response{Counter: RequestsCounter}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	lock.Unlock()

	return &jsonResponse, nil
}

func Counter(w http.ResponseWriter, r *http.Request) {
	response, err := UpdateCounter(w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(*response)
}

func main() {
	fmt.Println("Starting Server")

	exists, err := services.InitFile(logFilePath)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if exists {
		RequestsCounter, CounterID, err := services.ReadLogLines(logFilePath, windowSize)

		fmt.Printf("Actual counter %d\n", RequestsCounter)
		fmt.Printf("Counter ID is at %d\n", CounterID)

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		fmt.Println(RequestsCounter)
	}

	handler := http.HandlerFunc(Counter)
	http.Handle("/counter", handler)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server has encountered a problem")
	}
}

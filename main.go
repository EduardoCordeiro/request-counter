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

var logFilePath string = "requests.log"

func UpdateCounter(w http.ResponseWriter) (*[]byte, error) {
	lock.Lock()

	timestamp := time.Now().Local().Format(time.RFC3339)

	RequestsCounter, CounterID, err := services.ReadLogLines(logFilePath, windowSize)

	var logLine values.LogLine
	logLine.ID = CounterID
	logLine.Timestamp = timestamp

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

func startup() error {
	exists, err := services.InitFile(logFilePath)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if exists {
		RequestsCounter, _, err := services.ReadLogLines(logFilePath, windowSize)

		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Println(RequestsCounter)
	}

	return nil
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

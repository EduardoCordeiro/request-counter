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

var counter values.LogLine = values.LogLine{}
var lock sync.Mutex

var logFilePath string = "logs.txt"

func UpdateCounter(w http.ResponseWriter) (*[]byte, error) {
	lock.Lock()

	timestamp := time.Now().Local().Format(time.RFC3339)

	counter.Counter += 1
	counter.Timestamp = timestamp

	fmt.Println(counter)

	services.WriteToFile(logFilePath, &counter)

	jsonResponse, err := json.Marshal(counter)
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
		line, err := services.ReadLastLine(logFilePath)

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		err = json.Unmarshal([]byte(line), &counter)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			panic(err)
		}
		fmt.Println(line)
		fmt.Println(counter)
	}

	handler := http.HandlerFunc(Counter)
	http.Handle("/counter", handler)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server has encountered a problem")
	}
}

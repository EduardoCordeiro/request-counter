package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simpleinsurance/services"
	"simpleinsurance/values"
	"sync"
	"time"
)

const address string = "localhost:8080"

var newCounter values.LogLine = values.LogLine{}
var lock sync.Mutex

var logFilePath string = "logs.txt"

func WriteToFile() {

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write to file
	jsonData, err := json.Marshal(newCounter)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Convert the JSON data to a string.
	jsonString := string(jsonData)

	fmt.Println(jsonString)

	// Append the new line to the file.
	_, err = file.WriteString(jsonString + "\n")
	if err != nil {
		fmt.Printf("Error appending to file: %v\n", err)
		return
	}

	return
}

func UpdateCounter(w http.ResponseWriter) (*[]byte, error) {
	lock.Lock()

	timestamp := time.Now().Local().Format(time.RFC3339)

	newCounter.Counter += 1
	newCounter.Timestamp = timestamp

	fmt.Println(newCounter)

	WriteToFile()

	jsonResponse, err := json.Marshal(newCounter)
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

		err = json.Unmarshal([]byte(line), &newCounter)
		if err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			panic(err)
		}
		fmt.Println(line)
		fmt.Println(newCounter)
	}

	handler := http.HandlerFunc(Counter)
	http.Handle("/counter", handler)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server has encountered a problem")
	}
}

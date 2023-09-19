package main

import (
	"counter/handlers"
	"counter/services"
	"fmt"
	"log"
	"net/http"
)

const windowSize int = 60

var RequestsCounter int
var CounterID int

var logFilePath string = "requests.log"

func startup(path string) error {
	// Function to create or load a file, in case it exists
	err := services.CreateFile(path)
	if err != nil {
		log.Fatal(err)
		return err
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
	fmt.Println("Starting Server...")

	err := startup(logFilePath)
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

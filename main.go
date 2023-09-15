package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const address string = "localhost:8080"

var counter int = 0
var lock sync.Mutex

type Response struct {
	Counter int `json:"counter"`
}

func Counter(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	counter++
	response := Response{Counter: counter}
	lock.Unlock()

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	fmt.Print("Starting Server")

	handler := http.HandlerFunc(Counter)
	http.Handle("/counter", handler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Server has encountered a problem")
	}
}

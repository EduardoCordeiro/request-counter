package handlers

import (
	"counter/services"
	"counter/values"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var lock sync.Mutex

func UpdateCounter(w http.ResponseWriter, logFilePath string, windowSize int) (*[]byte, error) {
	lock.Lock()

	defer lock.Unlock()

	timestamp := time.Now().Local().Format(time.RFC3339)

	RequestsCounter, CounterID, err := services.ParseLogFile(logFilePath, windowSize)
	if err != nil {
		http.Error(w, "Error parsing the log file", http.StatusInternalServerError)
		return nil, err
	}

	var logLine values.LogLine
	logLine.ID = CounterID
	logLine.Timestamp = timestamp

	err = services.WriteToFile(logFilePath, &logLine)
	if err != nil {
		http.Error(w, "Error writing to the log file", http.StatusInternalServerError)
		return nil, err
	}
	// Create a Response value to output to the user
	response := values.Response{Counter: RequestsCounter}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	return &jsonResponse, nil
}

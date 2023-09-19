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
	/*
		Counter API Entrypoint

		This function is in charge of handling the locking and unlocking of the Mutex, used to control
		the acess to the file. Since the file contains the source of truth on the number of requests,
		we can only respond with the correct value if while we are parsing and writing to the file,
		no other request is doing the same.
		After obtaining the lock we parse the log file to get the number of request and the current id.
		We use the current id to write the new request into the file, and return the counter to the user.
	*/

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
	response := values.Response{Counter: RequestsCounter}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	return &jsonResponse, nil
}

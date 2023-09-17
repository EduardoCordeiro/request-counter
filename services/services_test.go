package services

import (
	"fmt"
	"os"
	"simpleinsurance/values"
	"testing"
	"time"
)

func generateLogs(times []int) []string {
	var logs []string
	var counter int = 0

	for _, time := range times {
		logLine := generateLogLine(counter, time)
		logs = append(logs, logLine)
		counter++
	}

	return logs
}

func generateLogLine(counter int, delta int) string {
	difference := time.Duration(time.Duration(delta) * time.Second)
	timestamp := time.Now().Add(-difference).Local().Format(time.RFC3339)

	logLine := fmt.Sprintf(`{"id":%d,"timestamp":"%s"}`, counter, timestamp)
	return logLine
}

func writeTestLogToFile(t *testing.T, filePath string, lines []string) {
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Error creating test log file: %v", err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Error writing test log line: %v", err)
		}
	}
}

func TestParseValidLogLine(t *testing.T) {
	// Test valid input
	validLogLine := `{"id":1,"timestamp":"2023-09-16T16:57:18+08:00"}`
	expectedLogLine := values.LogLine{
		ID:        1,
		Timestamp: "2023-09-16T16:57:18+08:00",
	}

	parsedLogLine, err := parseLogLine(validLogLine)
	if err != nil {
		t.Errorf("Error parsing valid log line: %v", err)
	}
	if parsedLogLine != expectedLogLine {
		t.Errorf("Expected: %+v, Got: %+v", expectedLogLine, parsedLogLine)
	}
}

func TestParseInvalidLogLine(t *testing.T) {
	// Test invalid input
	invalidLogLine := `invalid JSON`
	_, err := parseLogLine(invalidLogLine)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, but got nil")
	}
}

func TestReadLogLinesInside(t *testing.T) {
	// Create a temporary test log file with log lines
	logFilePath := "test.log"
	times := []int{20, 40}
	testLogLines := generateLogs(times)
	writeTestLogToFile(t, logFilePath, testLogLines)

	// Test reading log lines within a 1-minute window
	windowSize := 60
	counter, id, err := ReadLogLines(logFilePath, windowSize)
	if err != nil {
		t.Errorf("Error reading log lines: %v", err)
	}
	if counter != len(testLogLines) {
		t.Errorf("Expected counter: %d, Got: %d", len(testLogLines), counter)
	}
	if id != len(testLogLines) {
		t.Errorf("Expected ID: 2, Got: %d", id)
	}
}

func TestReadLogLinesOutside(t *testing.T) {
	// Create a temporary test log file with log lines
	logFilePath := "test.log"
	times := []int{20, 40}
	testLogLines := generateLogs(times)
	writeTestLogToFile(t, logFilePath, testLogLines)

	windowSize := 10
	counter, id, err := ReadLogLines(logFilePath, windowSize)

	if err != nil {
		t.Errorf("Error reading log lines: %v", err)
	}
	if counter != 0 {
		t.Errorf("Expected counter: 0, Got: %d", counter)
	}
	if id != 2 {
		t.Errorf("Expected ID: 2, Got: %d", id)
	}
}

func TestCheckDuration(t *testing.T) {

	// Test within the window
	line := generateLogLine(1, 40)
	windowSize := 300
	result, err := checkDuration(line, windowSize)
	fmt.Println("resultado")
	fmt.Println(result)
	if err != nil {
		t.Errorf("Error checking duration: %v", err)
	}
	if !result {
		t.Errorf("Expected result: true, Got: false")
	}

	// Test outside the window
	line = `{"id":1,"timestamp":"2023-09-16T16:56:00+08:00"}`
	windowSize = 20
	result, err = checkDuration(line, windowSize)
	if err != nil {
		t.Errorf("Error checking duration: %v", err)
	}
	if result {
		t.Errorf("Expected result: false, Got: true")
	}

	// Test invalid timestamp
	line = `{"id":1,"timestamp":"invalid-timestamp"}`
	windowSize = 60
	_, err = checkDuration(line, windowSize)
	if err == nil {
		t.Errorf("Expected error for invalid timestamp, but got nil")
	}
}

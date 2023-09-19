package services

import (
	"counter/values"
	"fmt"
	"os"
	"testing"
	"time"
)

type ParseLogLineTestCase struct {
	Description string
	Input       string
	Expected    values.LogLine
	ShouldError bool
}

type ReadLogLinesTestCase struct {
	Description     string
	LogLines        []string
	WindowSize      int
	ExpectedCounter int
	ExpectedID      int
}

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

func deleteLogFile(logFilePath string) error {
	err := os.Remove(logFilePath)
	if err != nil {
		return err
	}
	return nil
}

func TestParseLogLine(t *testing.T) {
	testCases := []ParseLogLineTestCase{
		{
			Description: "Valid input",
			Input:       `{"id":1,"timestamp":"2023-09-16T16:57:18+08:00"}`,
			Expected: values.LogLine{
				ID:        1,
				Timestamp: "2023-09-16T16:57:18+08:00",
			},
			ShouldError: false,
		},
		{
			Description: "Invalid input",
			Input:       `invalid`,
			ShouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			parsedLogLine, err := parseLogLine(testCase.Input)

			if testCase.ShouldError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if parsedLogLine != testCase.Expected {
					t.Errorf("Expected: %+v, Got: %+v", testCase.Expected, parsedLogLine)
				}
			}
		})
	}
}

func TestReadLogLines(t *testing.T) {
	testCases := []ReadLogLinesTestCase{
		{
			Description:     "Empty Log",
			LogLines:        []string{},
			WindowSize:      10,
			ExpectedCounter: 0,
			ExpectedID:      0,
		},
		{
			Description:     "Log Inside Window",
			LogLines:        generateLogs([]int{20, 40}),
			WindowSize:      60,
			ExpectedCounter: 2,
			ExpectedID:      2,
		},
		{
			Description:     "Log Outside Window",
			LogLines:        generateLogs([]int{20, 40}),
			WindowSize:      10,
			ExpectedCounter: 0,
			ExpectedID:      2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			logFilePath := "test.log"
			writeTestLogToFile(t, logFilePath, testCase.LogLines)

			counter, id, err := ParseLogFile(logFilePath, testCase.WindowSize)
			if err != nil {
				t.Errorf("Error reading log lines: %v", err)
			}
			if counter != testCase.ExpectedCounter {
				t.Errorf("Expected counter: %d, Got: %d", testCase.ExpectedCounter, counter)
			}
			if id != testCase.ExpectedID {
				t.Errorf("Expected ID: %d, Got: %d", testCase.ExpectedID, id)
			}

			err = deleteLogFile(logFilePath)
			if err != nil {
				t.Fatalf("Error removing the file: %v", err)
			}
		})
	}
}

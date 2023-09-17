package services

import (
	"fmt"
	"testing"
)

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

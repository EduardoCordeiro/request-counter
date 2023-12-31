package services

import (
	"fmt"
	"strings"
	"time"
)

func checkDuration(line string, windowSize int) (bool, error) {
	// Here we check if the current line we parsed is inside the window
	// Since we are saving the timestamp as RFC3339 (easier to read)
	// we convert those times and check if the duration between them is less
	// than the window size

	line = strings.ReplaceAll(line, "\x0A", "")
	line = strings.ReplaceAll(line, "\n", "")

	logLine, err := parseLogLine(line)
	if err != nil {
		return false, err
	}

	lineTime, err := time.Parse(time.RFC3339, logLine.Timestamp)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		return false, err
	}

	now := time.Now().Local().Format(time.RFC3339)
	nowTime, err := time.Parse(time.RFC3339, now)

	windowDuration := time.Duration(time.Duration(windowSize) * time.Second)

	duration := nowTime.Sub(lineTime)

	if duration <= windowDuration {
		return true, nil
	}

	return false, nil
}

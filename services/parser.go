package services

import (
	"counter/values"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func parseLogLine(line string) (values.LogLine, error) {
	var logLine values.LogLine
	err := json.Unmarshal([]byte(line), &logLine)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v\n", err)
		return values.LogLine{}, err
	}

	return logLine, nil
}

func ParseLogFile(filepath string, windowSize int) (int, int, error) {
	// This functin parses a file from the bottom up, to check which previous requests
	// are inside the window. It returns the number of requests inside the window, the total
	// request counter and an error.

	file, err := os.Open(filepath)

	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// Actual counter inside the window
	var counter int = 0
	// Last ID written
	var id int = -1
	line := ""
	var cursor int64 = 0
	stat, _ := file.Stat()
	fileSize := stat.Size()

	if fileSize == 0 {
		return 0, 0, nil
	}

	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		file.Read(char)

		// Stop when we find the end of the previous line
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			insideWindow, err := checkDuration(line, windowSize)
			if err != nil {
				return 0, 0, err
			}

			if id == -1 {
				logLine, err := parseLogLine(line)
				if err != nil {
					return 0, 0, err
				}

				id = logLine.ID + 1
			}

			// Only break when we find the first log after the window is over
			if insideWindow {
				counter++
				line = ""
			} else {
				break
			}
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -fileSize {
			break
		}
	}

	// Case where there is only one line
	// loop above does not reach because there is no new line before
	if line != "" {
		insideWindow, err := checkDuration(line, windowSize)
		if err != nil {
			return 0, 0, err
		}
		if insideWindow {
			counter++
		}

		if id == -1 {
			logLine, err := parseLogLine(line)
			if err != nil {
				return 0, 0, err
			}

			id = logLine.ID + 1
		}
	}

	return counter, id, nil
}

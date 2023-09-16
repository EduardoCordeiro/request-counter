package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"simpleinsurance/values"
	"strings"
	"time"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func PrintFileContents(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	// Print the file contents.
	fmt.Println("File Contents:")
	// Loop through each line in the file.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

func InitFile(filePath string) (bool, error) {
	fileExists := FileExists(filePath)

	var file *os.File
	var err error

	if fileExists {
		fmt.Println("File Exists!")
		file, err = os.Open(filePath)

		if err != nil {
			log.Fatal(err)
			return false, err
		}
		defer file.Close()

		return true, nil
	} else {
		fmt.Println("Creating a new data file!")
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("Error opening data file: %v\n", err)
			return false, err
		}
	}

	return false, nil
}

// this code will return the value of the last line, current counter for the ids of the requests
// Plus read the X amount of lines in the
func ReadLastLine(filepath string, windowSize int) (string, int, error) {

	// Actual counter inside the window
	var counter int = 0

	file, err := os.Open(filepath)

	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := file.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		file.Read(char)

		fmt.Printf("Char %s\n", string(char))

		// Stop when we find the end of the previous line
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {

			// Read lines, but only break when the Time.Duration > windowSize
			line = strings.ReplaceAll(line, "\x0A", "")
			line = strings.ReplaceAll(line, "\n", "")
			fmt.Println(line)

			var logLine values.LogLine
			err = json.Unmarshal([]byte(line), &logLine)
			if err != nil {
				fmt.Printf("Error unmarshaling JSON: %v\n", err)
				panic(err)
			}

			lineTime, err := time.Parse(time.RFC3339, logLine.Timestamp)

			if err != nil {
				fmt.Printf("Error parsing time: %v\n", err)
				return "", 0, err
			}

			now := time.Now().Local().Format(time.RFC3339)
			nowTime, err := time.Parse(time.RFC3339, now)

			windowDuration := time.Duration(time.Duration(windowSize) * time.Second)

			duration := nowTime.Sub(lineTime)

			fmt.Println("Duration")
			fmt.Println(duration)

			if duration <= windowDuration {
				counter++
				line = ""
			} else {
				break
			}
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize {
			break
		}
	}

	return line, counter, nil
}

func WriteToFile(filePath string, counter *values.LogLine) {

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write to file
	jsonData, err := json.Marshal(counter)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Convert the JSON data to a string.
	jsonString := string(jsonData)

	fmt.Println(jsonString)

	// Append the new line to the file.
	_, err = file.WriteString("\n" + jsonString)
	if err != nil {
		fmt.Printf("Error appending to file: %v\n", err)
		return
	}

	return
}

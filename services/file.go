package services

import (
	"counter/values"
	"encoding/json"
	"fmt"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateFile(filePath string) error {
	fileExists := FileExists(filePath)

	var err error

	if !fileExists {
		_, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("Error opening data file: %v\n", err)
			return err
		}
	}

	return nil
}

func WriteToFile(filePath string, counter *values.LogLine) error {
	// This functin writes a line to the log file, appending it to do the end of the file,
	// since we are parsing the file bottom-up

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return err
	}
	defer file.Close()

	// Write to file
	jsonData, err := json.Marshal(counter)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return err
	}

	// Convert the JSON data to a string.
	jsonString := string(jsonData)

	fmt.Println(jsonString)

	// Append the new line to the file.
	_, err = file.WriteString(jsonString + "\n")
	if err != nil {
		fmt.Printf("Error appending to file: %v\n", err)
		return err
	}

	return nil
}

func DeleteLogFile(logFilePath string) error {
	err := os.Remove(logFilePath)
	if err != nil {
		return err
	}
	return nil
}

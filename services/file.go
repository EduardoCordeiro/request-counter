package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"simpleinsurance/values"
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
		file, err = os.Open(filePath)

		if err != nil {
			log.Fatal(err)
			return false, err
		}
		defer file.Close()

		return true, nil
	} else {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("Error opening data file: %v\n", err)
			return false, err
		}
	}

	return false, nil
}

func WriteToFile(filePath string, counter *values.LogLine) error {

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

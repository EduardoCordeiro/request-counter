package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

func ReadLastLine(filepath string) (string, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return "", err
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

		// Stop when we find the end of the previous line
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize {
			break
		}
	}

	return line, nil
}

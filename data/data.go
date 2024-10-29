package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	filename = "output.csv"
	Header   = []string{"ID", "description", "createdAt", "isCompleted"}
)

func openTasksFile() (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error while trying to open tasks file")
		return nil, err
	}

	return f, nil
}

func writeToTasksFile(writer *csv.Writer, content []string) error {
	f, openErr := openTasksFile()

	if openErr != nil {
		return openErr
	}

	defer f.Close()

	writeErr := writer.Write(content)

	if writeErr != nil {
		return writeErr
	}

	return nil
}

func ReadTasks() ([][]string, *os.File, error) {
	f, _ := openTasksFile()

	defer f.Close()

	bytesData, readErr := io.ReadAll(f)
	if readErr != nil {
		fmt.Println("Error while trying to read file")

		return nil, nil, readErr
	}

	reader := csv.NewReader(bytes.NewReader(bytesData))

	var tasks [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading CSV data:", err)
			return nil, nil, err
		}

		tasks = append(tasks, record)
	}

	// The file is just created or empty, add headers
	if len(tasks) == 0 {
		writer := csv.NewWriter(f)

		tasks = append(tasks, Header)
		writeToTasksFile(writer, Header)

		writer.Flush()
	}

	// Only return tasks, without the header
	return tasks[1:], f, nil

}

func CreateTask(description string) {
	tasks, f, _ := ReadTasks()

	// Closing the file before opening it again, to avoid concurrent writing
	f.Close()

	newTask := []string{strconv.FormatInt(int64(len(tasks))+1, 10), description, time.Now().Format(time.RFC3339), "false"}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error while opening tasks file for write")
	}

	defer f.Close()

	w := csv.NewWriter(f)

	w.Write(newTask)

	w.Flush()
}

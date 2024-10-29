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

	// Insert the new task at the beginning of the slice
	newTask := []string{strconv.FormatInt(int64(len(tasks))+1, 10), description, time.Now().Format(time.RFC3339), "false"}
	tasks = append([][]string{newTask}, tasks...)

	// Reopen the file in truncate mode to overwrite the contents
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error while opening tasks file for write")
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// Write header, followed by all tasks with the latest task first
	if err := w.Write(Header); err != nil {
		fmt.Println("Error writing header:", err)
		return
	}
	if err := w.WriteAll(tasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}
	w.Flush()

	fmt.Println("Task created")
}

func CompleteTask(taskId int) error {
	tasks, f, _ := ReadTasks()
	defer f.Close()

	// Open the file in truncate mode to rewrite all tasks
	f, _ = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()

	for i, task := range tasks {
		// Parse the ID field to an integer
		id, err := strconv.Atoi(task[0])
		if err != nil {
			return fmt.Errorf("invalid task ID in file: %v", err)
		}

		// Check if the task ID matches the input taskId
		if id == taskId {
			// Mark the task as completed
			tasks[i][3] = "true"
			break
		}
	}

	// Write updated tasks back to the file, including the header
	w := csv.NewWriter(f)
	if err := w.Write(Header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}
	if err := w.WriteAll(tasks); err != nil {
		return fmt.Errorf("error writing tasks: %v", err)
	}
	w.Flush()

	// fmt.Printf("Task with ID %d marked as completed.\n", taskId)
	return nil
}

func DeleteTask(taskId int) error {
	// Read all tasks
	tasks, _, err := ReadTasks()
	if err != nil {
		return fmt.Errorf("error reading tasks: %v", err)
	}

	// Find the index of the task with the specified ID
	indexToDelete := -1
	for i, task := range tasks {
		id, err := strconv.Atoi(task[0]) // Convert ID from string to int
		if err != nil {
			return fmt.Errorf("invalid task ID in file: %v", err)
		}

		if id == taskId {
			indexToDelete = i
			break
		}
	}

	// If the task was not found, return an error
	if indexToDelete == -1 {
		return fmt.Errorf("task with ID %d not found", taskId)
	}

	// Remove the task from the slice
	tasks = append(tasks[:indexToDelete], tasks[indexToDelete+1:]...)

	// Open the file in truncate mode to rewrite the tasks
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %v", err)
	}
	defer f.Close()

	// Write the header and remaining tasks back to the file
	w := csv.NewWriter(f)
	if err := w.Write(Header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}
	if err := w.WriteAll(tasks); err != nil {
		return fmt.Errorf("error writing tasks: %v", err)
	}
	w.Flush()

	// fmt.Printf("Task with ID %d deleted successfully.\n", taskId)
	return nil
}

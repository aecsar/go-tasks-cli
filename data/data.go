package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var (
	header = []string{"ID", "description", "createdAt", "isCompleted"}
)

func CreateTask(description string) {
	f, err := os.Create("output.csv")

	if err != nil {
		fmt.Println("Error while creating file", err)
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)

	w.Write(header)

	w.Write([]string{"1", description, time.Now().String(), "false"})

	w.Flush()
	// writeError :=
	// if writeError != nil {
	// 	fmt.Println("Error writing record to CSV:", err)
	// }

	fmt.Println("created task : ", description)
}

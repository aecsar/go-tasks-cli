/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"aecsar/tasks/data"
	"fmt"
	"os"
	"time"

	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/mergestat/timediff"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get tasks list",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("list called")

		tasks, f, _ := data.ReadTasks()
		defer f.Close()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

		header := []string{"ID", "Task", "Creation", "Done"}

		tasks = append([][]string{header}, tasks...)

		for rowIdx, task := range tasks {
			var line string

			for propIndex, prop := range task {
				if propIndex == 2 && rowIdx > 0 {
					parsedTime, _ := time.Parse(time.RFC3339, prop)

					// if parseErr != nil {
					// 	fmt.Println("Error parsing time : ", parseErr)
					// }

					prop = timediff.TimeDiff(parsedTime)
				}

				line += prop + "\t"
			}

			fmt.Fprintln(w, line)
		}

		w.Flush()

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

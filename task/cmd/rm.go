package cmd

import (
	"fmt"
	"github.com/abhyuditjain/gophercices/task/db"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Deletes tasks",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument: ", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number: ", id)
				continue
			}
			task := tasks[id-1]
			err := db.DeleteTask(task.Key, nil)
			if err != nil {
				fmt.Printf("Failed to delete \"%s\". Error: %s\n", task.Value, err)
			} else {
				fmt.Printf("You have deleted the \"%s\" task.\n", task.Value)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}

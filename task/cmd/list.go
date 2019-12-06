package cmd

import (
	"fmt"
	"github.com/abhyuditjain/gophercices/task/db"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the TODOs",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete! Go for a vacation 🏖️")
			return
		}
		fmt.Println("You have the following tasks:")
		for i, v := range tasks {
			fmt.Printf("%d. %s\n", i+1, v.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

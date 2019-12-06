package cmd

import (
	"fmt"
	"github.com/abhyuditjain/gophercices/task/db"
	"github.com/spf13/cobra"
	"os"
)

// completedCommand represents the completed command
var completedCommand = &cobra.Command{
	Use:   "completed",
	Short: "Lists the completed TODOs for today",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.CompletedTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have completed 0 tasks today")
			return
		}
		fmt.Println("You have finished the following tasks today:")
		for _, v := range tasks {
			fmt.Printf("- %s\n", v.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(completedCommand)
}

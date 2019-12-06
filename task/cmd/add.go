package cmd

import (
	"fmt"
	"github.com/abhyuditjain/gophercices/task/db"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a TODO to the list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

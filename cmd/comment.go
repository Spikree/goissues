package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Spikree/git-issues/helper"
	"github.com/Spikree/git-issues/types"
	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment <id> text",
	Short: "add a comment to an issue",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, text := args[0], args[1]

		db, err := helper.LoadDB()
		if err != nil {
			fmt.Println("Database connection error:", err)
			os.Exit(1)
		}

		parentIssue, exists := db[id]

		newComment := types.Comment{
			Author:    helper.GitUser(),
			Text:      text,
			CreatedAt: time.Now(),
		}

		if exists {
			parentIssue.Comments = append(parentIssue.Comments, newComment)
			db[id] = parentIssue
		} else {
		searchSubIssues:
			for parentId, parent := range db {
				for i := range parent.SubIssues {
					if parent.SubIssues[i].ID != id {
						continue
					}

					parent.SubIssues[i].Comments = append(
						parent.SubIssues[i].Comments,
						newComment,
					)
					db[parentId] = parent
					exists = true
					break searchSubIssues
				}
			}
		}

		if !exists {
			fmt.Println("issue not found")
			os.Exit(1)
		}

		if err := helper.SaveDB(db); err != nil {
			fmt.Println("Error adding a comment", err)
			os.Exit(1)
		}

		fmt.Println("comment added to", id)
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

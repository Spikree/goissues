/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Spikree/git-issues/helper"
	"github.com/spf13/cobra"
)

// unresolveCmd represents the unresolve command
var unresolveCmd = &cobra.Command{
	Use:   "unresolve",
	Short: "Use this command to unresolve an issue in your project",
	Long:  `use this command to unresolve an issue in your project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		db, err := helper.LoadDB()
		if err != nil {
			fmt.Println("Database connection error:", err)
			os.Exit(1)
		}

		parentIssue, exists := db[id]

		if exists {
			if parentIssue.Status == "Open" {
				fmt.Println("Issue is already open", id)
				return
			}

			if parentIssue.Status != "Resolved" {
				fmt.Fprintln(os.Stderr, "Cannot reopen issue with status:", parentIssue.Status)
				os.Exit(1)
			}

			parentIssue.Status = "Open"
			db[id] = parentIssue
		} else {
		searchSubIssues:
			for parentID, parent := range db {
				for i := range parent.SubIssues {
					// Update the slice element, not a copy or its parent.
					sub := &parent.SubIssues[i]
					if sub.ID != id {
						continue
					}

					if sub.Status == "Open" {
						fmt.Println("Issue is already open", id)
						return
					}

					if sub.Status != "Resolved" {
						fmt.Fprintln(os.Stderr, "Cannot open issue with status:", sub.Status)
						os.Exit(1)
					}

					sub.Status = "Open"
					db[parentID] = parent

					exists = true
					// Leave both loops and continue to the shared save below.
					break searchSubIssues
				}
			}
		}

		if !exists {
			fmt.Fprintln(os.Stderr, "Issue not found", id)
			os.Exit(1)
		}

		if err := helper.SaveDB(db); err != nil {
			fmt.Fprintln(os.Stderr, "Error saving issue", err)
			os.Exit(1)
		}

		fmt.Println("Issue Opened", id)
	},
}

func init() {
	rootCmd.AddCommand(unresolveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unresolveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unresolveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

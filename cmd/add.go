/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/Spikree/git-issues/helper"
	"github.com/Spikree/git-issues/types"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Use this command to add an issue to your project",
	Long: `use this command to add an issue to your project use "git issues add ...", this
	command takes in a title and a body`,
	Run: func(cmd *cobra.Command, args []string) {
		// here we get the title and body from the command line
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")

		// Load all the data into ram
		db, err := helper.LoadDB()
		if err != nil {
			fmt.Println("Database connection error:", err)
			os.Exit(1)
		}

		// Generate a unique primary key
		b := make([]byte, 3)
		_, _ = rand.Read(b)
		issueId := hex.EncodeToString(b)

		newIssue := types.Issue{
			ID:        issueId,
			Title:     title,
			Body:      body,
			Status:    "Open",
			CreatedAt: time.Now(),
			SubIssues: []types.SubIssue{},
			Comments:  []types.Comment{},
		}

		db[issueId] = newIssue

		if err := helper.SaveDB(db); err != nil {
			fmt.Println("Failed to save transaction records:", err)
			os.Exit(1)
		}

		fmt.Println("Issue created: ", issueId)
	},
}

var addSubIssueCmd = &cobra.Command{
	Use:   "sub-issue",
	Short: "Add a sub issue to the main issue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// parent issue id for the sub issue
		id := args[0]
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")

		db, err := helper.LoadDB()
		if err != nil {
			fmt.Println("Database connection error", err)
			os.Exit(1)
		}

		parentIssue, exists := db[id]
		if !exists {
			fmt.Println("parent issue not found")
			os.Exit(1)
		}

		b := make([]byte, 3)
		_, _ = rand.Read(b)
		subIssueId := hex.EncodeToString(b)

		newSubIssue := types.SubIssue{
			ID:        subIssueId,
			Title:     title,
			Body:      body,
			Status:    "Open",
			CreatedAt: time.Now(),
			Comments:  []types.Comment{},
		}

		parentIssue.SubIssues = append(parentIssue.SubIssues, newSubIssue)

		db[id] = parentIssue

		if err := helper.SaveDB(db); err != nil {
			fmt.Println("Error creating sub issue", err)
			os.Exit(1)
		}

		fmt.Println("sub issue created", subIssueId)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("title", "t", "", "issue title")
	addCmd.Flags().StringP("body", "b", "", "issue description")
	addCmd.MarkFlagRequired("title")

	addSubIssueCmd.Flags().StringP("title", "t", "", "sub issue title")
	addSubIssueCmd.Flags().StringP("body", "b", "", "sub issue description")
	addSubIssueCmd.MarkFlagRequired("title")
	addCmd.AddCommand(addSubIssueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

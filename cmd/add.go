/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type Comment struct {
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type SubIssue struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	Comments  []Comment `json:"comments,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Issue struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Status    string     `json:"Status"`
	CreatedAt time.Time  `json:"created_at"`
	SubIssues []SubIssue `json:"sub_issues,omitempty"`
	Comments  []Comment  `json:"comments,omitempty"`
}

// Database scheme config
const DbPath = ".issues_db.json"

// Load database state from disk into memory
func loadDB() (map[string]Issue, error) {
	db := make(map[string]Issue)

	// if the database isnt initilised yet, return initialised map
	if _, err := os.Stat(DbPath); os.IsNotExist(err) {
		return db, nil
	}

	data, err := os.ReadFile(DbPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &db)
	return db, err
}

func saveDB(db map[string]Issue) error {
	// marshalIndent keeps the file human readable in Git diffs
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(DbPath, data, 0o644)
}

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
		db, err := loadDB()
		if err != nil {
			fmt.Println("Database connection error:", err)
			os.Exit(1)
		}

		// Generate a unique primary key
		b := make([]byte, 3)
		_, _ = rand.Read(b)
		issueId := hex.EncodeToString(b)

		newIssue := Issue{
			ID:        issueId,
			Title:     title,
			Body:      body,
			Status:    "Open",
			CreatedAt: time.Now(),
			SubIssues: []SubIssue{},
			Comments:  []Comment{},
		}

		db[issueId] = newIssue

		if err := saveDB(db); err != nil {
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

		path := filepath.Join(".issues", id+".md")

		data, err := os.ReadFile(path)

		if err != nil {
			fmt.Println("issue not found", id)
			os.Exit(1)
		}

		b := make([]byte, 3)
		_, _ = rand.Read(b)
		subIssueId := hex.EncodeToString(b)

		entry := fmt.Sprintf("\n### Sub Issue: %s\n- **ID:** %s\n- **Status:** open\n- **Created:** %s\n\n%s\n",
			title, subIssueId, time.Now().Format("2006-01-02"), body)

		var finalContent []byte

		subIssueIdx := bytes.Index(data, []byte("\n### Sub Issue:"))

		insertIdx := -1

		if subIssueIdx != -1 {
			insertIdx = subIssueIdx
		} else {
			insertIdx = subIssueIdx
		}

		if insertIdx != -1 {
			finalContent = append(finalContent, data[:insertIdx]...)
			finalContent = append(finalContent, []byte(entry)...)
			finalContent = append(finalContent, data[insertIdx:]...)
		} else {
			finalContent = append(data, []byte(entry)...)
		}

		if err := os.WriteFile(path, finalContent, 0o644); err != nil {
			fmt.Println("error:", err)
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

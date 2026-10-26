/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

		// here we create a dir for all the issues and see if we get any errors and if we do we use os.exit(1)
		if err := os.MkdirAll(".issues", 0o755); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		// here we create a len 6 hex code
		b := make([]byte, 3)
		rand.Read(b)
		id := hex.EncodeToString(b)

		// here we store everything in a exact typed format
		content := fmt.Sprintf("### Issue: %s\n- **ID:** %s\n- **Status:** open\n- **Created:** %s\n\n%s\n",
			title, id, time.Now().Format("2006-01-02"), body)

		// here we write to the path
		path := filepath.Join(".issues", id+".md")

		// we check for any errors while writing
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("issue Created", path)
		// fmt.Println("Adding issue:", title)
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

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("Issue not found", id)
			os.Exit(1)
		}

		b := make([]byte, 3)
		_, _ = rand.Read(b)
		subIssueId := hex.EncodeToString(b)

		content := fmt.Sprintf("\n### Sub Issue: %s\n- **ID:** %s\n- **Status:** open\n- **Created:** %s\n\n%s\n",
			title, subIssueId, time.Now().Format("2006-01-02"), body)

		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Println("error opening file:", err)
			os.Exit(1)
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
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

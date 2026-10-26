package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Spikree/git-issues/helper"
	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment <id> text",
	Short: "add a comment to an issue",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id, text := args[0], args[1]

		path := filepath.Join(".issues", id+".md")

		data, err := os.ReadFile(path)

		if err != nil {
			fmt.Println("issue not found", id)
			os.Exit(1)
		}

		// we get the username from the helper function called gitUser
		entry := fmt.Sprintf("\n## Comment by %s on %s\n\n%s\n",
			helper.GitUser(), time.Now().Format("2006-01-02 15:04"), text)

		var finalContent []byte

		commentIdx := bytes.Index(data, []byte("\n## Comment by"))
		subIssueIdx := bytes.Index(data, []byte("\n### Sub Issue:"))

		insertIdx := -1

		if commentIdx != -1 {
			insertIdx = commentIdx
		} else if subIssueIdx != -1 {
			insertIdx = subIssueIdx
		}

		if insertIdx != -1 {
			// Split and inject the new comment entry at the calculated index position
			finalContent = append(finalContent, data[:insertIdx]...)
			finalContent = append(finalContent, []byte(entry)...)
			finalContent = append(finalContent, data[insertIdx:]...)
		} else {
			// fallback append if no comments or sub-issues exist yet
			finalContent = append(data, []byte(entry)...)
		}

		if err := os.WriteFile(path, finalContent, 0o644); err != nil {
			fmt.Println("error", err)
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

package cmd

import (
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

		if _, err := os.Stat(path); err != nil {
			fmt.Println("issue not found:", id)
			os.Exit(1)
		}

		// we use the path to append the comment at the end of the file, and we use O_WRONGLY so we dont touch the existing
		// issue structure
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)

		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}

		// we get the username from the helper function called gitUser
		entry := fmt.Sprintf("\n## Comment by %s on %s\n\n%s\n",
			helper.GitUser(), time.Now().Format("2006-01-02 15:04"), text)

		// we use the writestring function to add the comment
		if _, err := f.WriteString(entry); err != nil {
			f.Close()
			fmt.Println("error", err)
			os.Exit(1)
		}

		f.Close()
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

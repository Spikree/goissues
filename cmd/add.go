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
	Long:  `use this command to add an issue to your project use "git issues add ..."`,
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
		content := fmt.Sprintf(`---
		id: %s
		title: %s
		status: open
		created: %s
		---
		%s
		`, id, title, time.Now().Format("2006-01-02"), body)

		// here we write to the path
		path := filepath.Join(".issues", id+".md")

		// we check for any errors while writing
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Created", path)

		// fmt.Println("Adding issue:", title)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("title", "t", "", "issue title")
	addCmd.Flags().StringP("body", "b", "", "issue description")
	addCmd.MarkFlagRequired("title")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve <id>",
	Short: "use this to resolve any issue in your project",
	Long:  `use this command to resolve any ongoing issue in your project`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		path := filepath.Join(".issues", id+".md")

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("issue not found", id)
			os.Exit(1)
		}

		text := string(data)
		if !strings.Contains(text, "status: open") {
			fmt.Println("issue is not open", id)
			os.Exit(1)
		}

		text = strings.Replace(text, "status: open", "status: resolved", 1)

		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Println("Resolved", id)
	},
}

func init() {
	rootCmd.AddCommand(resolveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resolveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resolveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

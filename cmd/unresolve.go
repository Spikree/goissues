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

// unresolveCmd represents the unresolve command
var unresolveCmd = &cobra.Command{
	Use:   "unresolve",
	Short: "Use this command to unresolve an issue in your project",
	Long:  `use this command to unresolve an issue in your project`,
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
		if !strings.Contains(text, "status: resolved") {
			fmt.Println("issue is not resolved", id)
			os.Exit(1)
		}

		text = strings.Replace(text, "status: resolved", "status: open", 1)

		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}
		fmt.Println("issue opened", id)
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

/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/Spikree/git-issues/helper"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your issues",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := helper.LoadDB()
		if err != nil {
			fmt.Println("Database connection error", err)
			os.Exit(1)
		}

		if len(db) == 0 {
			fmt.Println("No issues found")
			return
		}

		ids := make([]string, 0, len(db))
		for id := range db {
			ids = append(ids, id)
		}
		sort.Strings(ids)

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tTITLE")
		fmt.Fprintln(w, "──\t──────\t─────")
		for _, id := range ids {
			issue := db[id]
			fmt.Fprintf(w, "%s\t%s\t%s\n", listCell(id), listCell(issue.Status), listCell(issue.Title))

			for i, sub := range issue.SubIssues {
				branch := "  ├─ "
				if i == len(issue.SubIssues)-1 {
					branch = "  └─ "
				}
				fmt.Fprintf(w, "%s%s\t%s\t%s\n", branch, listCell(sub.ID), listCell(sub.Status), listCell(sub.Title))
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "Error writing issue list:", err)
			os.Exit(1)
		}
	},
}

// Keep stored text on one line and strip terminal control characters.
func listCell(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\t' || r == '\n' || r == '\r':
			return ' '
		case r < ' ' || (r >= '\x7f' && r <= '\x9f'):
			return -1
		default:
			return r
		}
	}, s)
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

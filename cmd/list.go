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
	"github.com/Spikree/git-issues/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "List all your issues",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		status := "all"
		if len(args) > 0 {
			status = strings.ToLower(args[0])
		}
		switch status {
		case "all", "open", "resolved":
		default:
			fmt.Fprintf(cmd.ErrOrStderr(), "Unknown status %q: use all, open, or resolved\n", status)
			os.Exit(1)
		}
		matchesStatus := func(issueStatus string) bool {
			return status == "all" || strings.EqualFold(issueStatus, status)
		}

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

			// store all the visible sub issues here
			var visibleSubs []types.SubIssue
			// go through the sub issues to see if their status match the args status if they do store them in the
			// visibleSubs
			for _, sub := range issue.SubIssues {
				if matchesStatus(sub.Status) {
					visibleSubs = append(visibleSubs, sub)
				}
			}
			// if the main issue doesnt match and its sub issues also dont match status dont do anything
			if !matchesStatus(issue.Status) && len(visibleSubs) == 0 {
				continue
			}

			// else we print the issue with matching status and the issue which might not have a mathcing status but has
			// sub issues whos status match
			fmt.Fprintf(w, "%s\t%s\t%s\n", listCell(id), listCell(issue.Status), listCell(issue.Title))
			for i, sub := range visibleSubs {
				branch := "  ├─ "
				if i == len(visibleSubs)-1 {
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

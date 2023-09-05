package cmd

import (
	"strings"

	"github.com/2g4/console-ai/data"
	"github.com/spf13/cobra"
)

// qCmd represents the q command
var qCmd = &cobra.Command{
	Use:   "q",
	Short: "Start conversation with AI",
	Long:  "Start conversation with AI",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 && args[0] != "" {
			// If question is passed as argument
			// Run the query
			data.RunQuery(strings.Join(args, " "))
		} else {
			// If no question is passed as argument
			// Read question from user
			data.ReadQuestion()
		}
	},
}

func init() {

	rootCmd.AddCommand(qCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

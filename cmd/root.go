package cmd

import (
	"os"

	"github.com/2g4/console-ai/data"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "console-ai",
	Short: "Console tool to interact with ChatGPT",
	Long:  `Console Tool to interact with ChatGPT is a command line interface that allows users to communicate with the ChatGPT natural language processing (NLP) engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		data.ReadQuestion()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// InitIfRequired checks if settings table is created and if not,
// it creates it and populates with user & default values
func InitIfRequired() {
	if data.IsSettingsTableCreated() {
		data.FetchSettings()
		return
	}
	data.InitOrDie()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

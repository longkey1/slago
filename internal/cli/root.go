package cli

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	token   string
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "slago",
		Short: "Slack Log Collector CLI",
		Long: `slago is a CLI tool for collecting Slack messages.

It supports:
  - Fetching messages by URL
  - Collecting messages for date ranges
  - Filtering by author and mentions
  - Thread expansion`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.slago.yaml)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Slack API token (overrides SLACK_API_TOKEN)")

	// Add subcommands
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}

// Execute runs the root command
func Execute() error {
	return NewRootCmd().Execute()
}

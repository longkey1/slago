package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (set by ldflags)
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// SetVersion sets the version information
func SetVersion(version, commit, date string) {
	if version != "" {
		Version = version
	}
	if commit != "" {
		Commit = commit
	}
	if date != "" {
		Date = date
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("slago version %s\n", Version)
			fmt.Printf("  commit: %s\n", Commit)
			fmt.Printf("  built:  %s\n", Date)
		},
	}
}

package main

import (
	"os"

	"github.com/longkey1/slago/cmd"
)

// Version information (set by ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

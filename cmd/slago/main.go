package main

import (
	"os"

	"github.com/longkey1/slago/internal/cli"
)

// Version information (set by ldflags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.SetVersion(version, commit, date)

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

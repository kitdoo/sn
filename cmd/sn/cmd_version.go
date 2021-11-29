package main

import (
	"fmt"

	"github.com/kitdoo/sn/internal/version"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run:   cmdVersionRun,
}

func init() {
	RootCmd.AddCommand(cmdVersion)
}

func cmdVersionRun(_ *cobra.Command, _ []string) {
	fmt.Printf("%s - %s\n", version.AppName, version.FullVersion())
}

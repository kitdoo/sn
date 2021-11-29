package main

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sn",
	Short: "sn",
}

func init() {
	RootCmd.SuggestionsMinimumDistance = 1
	RootCmd.SilenceUsage = true
	RootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file or directory with config files")
}

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

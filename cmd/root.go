package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:          "smts",
	Short:        "Sign Me This Shit",
	Long:         "SMTS is a tool to generate and sign attendance sheets for FIPs students.",
	SilenceUsage: true,
}

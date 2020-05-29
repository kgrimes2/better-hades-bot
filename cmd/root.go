package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "bhb",
		Short: "Better Hades Bot",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ikun",
	Short: "Ikun Messenger is a distributed IM system, based on Raft",
	Long:  "Ikun Messenger is a distributed IM system, based on Raft",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

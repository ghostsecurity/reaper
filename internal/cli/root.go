package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/version"
)

var rootCmd = &cobra.Command{
	Use:   "reaper",
	Short: "Ghost Security MITM HTTPS proxy tool for application security testing",
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("reaper %s\n", version.Version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

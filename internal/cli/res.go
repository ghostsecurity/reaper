package cli

import (
	"github.com/spf13/cobra"
)

var resCmd = &cobra.Command{
	Use:   "res <id>",
	Short: "Show raw HTTP response for an entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runRes,
}

func init() {
	rootCmd.AddCommand(resCmd)
}

func runRes(cmd *cobra.Command, args []string) error {
	return fetchAndPrint(args[0], "res")
}

package cli

import (
	"github.com/spf13/cobra"
)

var reqCmd = &cobra.Command{
	Use:   "req <id>",
	Short: "Show raw HTTP request for an entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runReq,
}

func init() {
	rootCmd.AddCommand(reqCmd)
}

func runReq(cmd *cobra.Command, args []string) error {
	return fetchAndPrint(args[0], "req")
}

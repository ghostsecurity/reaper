package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running daemon",
	RunE:  runShutdown,
}

var shutdownCmd = &cobra.Command{
	Use:    "shutdown",
	Short:  "Stop the running daemon",
	Hidden: true,
	RunE:   runShutdown,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(shutdownCmd)
}

func runShutdown(cmd *cobra.Command, args []string) error {
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	client := daemon.NewClient(dataDir)
	resp, err := client.Send(daemon.Request{Command: "shutdown"})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}

	if !resp.OK {
		return fmt.Errorf("shutdown failed: %s", resp.Error)
	}

	fmt.Println("reaper daemon stopped")
	return nil
}

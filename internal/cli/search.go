package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ghostsecurity/reaper/internal/daemon"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search proxy log entries",
	RunE:  runSearch,
}

var (
	searchMethod  string
	searchHost    string
	searchDomains []string
	searchPath    string
	searchStatus  int
	searchLimit   int
)

func init() {
	searchCmd.Flags().StringVar(&searchMethod, "method", "", "Filter by HTTP method")
	searchCmd.Flags().StringVar(&searchHost, "host", "", "Filter by host (supports * wildcard)")
	searchCmd.Flags().StringSliceVar(&searchDomains, "domains", nil, "Filter by domain suffix")
	searchCmd.Flags().StringVar(&searchPath, "path", "", "Filter by path prefix or glob")
	searchCmd.Flags().IntVar(&searchStatus, "status", 0, "Filter by status code")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 100, "Max results")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	dataDir, err := daemon.DataDir()
	if err != nil {
		return err
	}

	params, _ := json.Marshal(daemon.SearchRequestParams{
		Method:  searchMethod,
		Host:    searchHost,
		Domains: searchDomains,
		Path:    searchPath,
		Status:  searchStatus,
		Limit:   searchLimit,
	})

	client := daemon.NewClient(dataDir)
	resp, err := client.Send(daemon.Request{Command: "search", Params: params})
	if err != nil {
		return fmt.Errorf("no running daemon found: %w", err)
	}
	if !resp.OK {
		return fmt.Errorf("%s", resp.Error)
	}

	var entries []entryRow
	if err := json.Unmarshal(resp.Data, &entries); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	printTable(entries)
	return nil
}

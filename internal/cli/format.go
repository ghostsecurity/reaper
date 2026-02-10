package cli

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
)

// entryRow is a subset of storage.Entry used for table display (deserialized from JSON).
type entryRow struct {
	ID         int64       `json:"ID"`
	Method     string      `json:"Method"`
	Scheme     string      `json:"Scheme"`
	Host       string      `json:"Host"`
	Path       string      `json:"Path"`
	Query      string      `json:"Query"`
	StatusCode int         `json:"StatusCode"`
	DurationMs int64       `json:"DurationMs"`
	// These fields are present but not used for table display
	RequestHeaders  http.Header `json:"RequestHeaders"`
	RequestBody     []byte      `json:"RequestBody"`
	ResponseHeaders http.Header `json:"ResponseHeaders"`
	ResponseBody    []byte      `json:"ResponseBody"`
}

// entryFull has all fields for raw display.
type entryFull = entryRow

func printTable(entries []entryRow) {
	if len(entries) == 0 {
		fmt.Println("no entries found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
		"ID", "METHOD", "HOST", "PATH", "STATUS",
		pad("MS", 6), pad("REQ", 7), pad("RES", 7))

	for _, e := range entries {
		path := e.Path
		if e.Query != "" {
			path += "?" + e.Query
		}
		if len(path) > 60 {
			path = path[:57] + "..."
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%6d\t%7s\t%7s\t\n",
			e.ID, e.Method, e.Host, path, e.StatusCode, e.DurationMs,
			formatSize(len(e.RequestBody)), formatSize(len(e.ResponseBody)))
	}

	w.Flush()
	fmt.Printf("\n%d entries\n", len(entries))
}

func printRawRequest(e entryFull) {
	path := e.Path
	if e.Query != "" {
		path += "?" + e.Query
	}
	if path == "" {
		path = "/"
	}

	fmt.Printf("%s %s HTTP/1.1\r\n", e.Method, path)
	fmt.Printf("Host: %s\r\n", e.Host)
	printHeaders(e.RequestHeaders)
	fmt.Print("\r\n")
	if len(e.RequestBody) > 0 {
		fmt.Printf("%s", e.RequestBody)
	}
}

func printRawResponse(e entryFull) {
	statusText := http.StatusText(e.StatusCode)
	if statusText == "" {
		statusText = "Unknown"
	}
	fmt.Printf("HTTP/1.1 %d %s\r\n", e.StatusCode, statusText)
	printHeaders(e.ResponseHeaders)
	fmt.Print("\r\n")
	if len(e.ResponseBody) > 0 {
		fmt.Printf("%s", e.ResponseBody)
	}
}

func pad(s string, width int) string {
	return fmt.Sprintf("%*s", width, s)
}

func formatSize(n int) string {
	switch {
	case n == 0:
		return "-"
	case n < 1024:
		return fmt.Sprintf("%dB", n)
	case n < 1024*1024:
		return fmt.Sprintf("%.1fKB", float64(n)/1024)
	default:
		return fmt.Sprintf("%.1fMB", float64(n)/(1024*1024))
	}
}

func printHeaders(h http.Header) {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range h[k] {
			fmt.Printf("%s: %s\r\n", k, v)
		}
	}
}

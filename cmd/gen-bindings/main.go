package main

import (
	"fmt"
	"os"

	"github.com/ghostsecurity/reaper/backend/server/bindings"
)

func main() {
	if err := bindings.Generate(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

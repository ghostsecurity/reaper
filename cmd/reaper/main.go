package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ghostsecurity/reaper/backend/server"
	"github.com/ghostsecurity/reaper/frontend"

	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/settings"
)

func main() {

	userSettings, err := settings.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading settings: %s\n", err)
		os.Exit(1)
	}

	logger := log.New(os.Stderr)
	level := userSettings.Get().LogLevel

	if logLevelName := os.Getenv("REAPER_LOG_LEVEL"); logLevelName != "" {
		level = log.ParseLevel(logLevelName)
	}

	logger.SetLevel(level)
	//logger.SetLevel(log.LevelDebug)

	logger.Info("User settings loaded...")
	logger.Printf(level, "Log level is set (see left)")

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv := server.New(ctx, frontend.Static, logger, userSettings)

	addr := "localhost:31337"
	printAddr := addr
	if _, err := os.Stat("/.dockerenv"); err == nil {
		addr = ":31337"
	}

	message([]string{
		"Welcome to Reaper!",
		"",
		"GUI running at http://" + printAddr,
	})

	if err := srv.Start(addr); err != nil {
		if !errors.Is(err, context.Canceled) {
			_, _ = fmt.Fprintf(os.Stderr, "Error starting server: %s\n", err)
			os.Exit(1)
		}
	}

	logger.Info("Exited cleanly.")
}

func message(lines []string) {
	w := 60
	fmt.Println("┌" + strings.Repeat("─", w-2) + "┐")
	for _, line := range lines {
		fmt.Println("│ " + line + strings.Repeat(" ", w-4-len(line)) + " │")
	}
	fmt.Println("└" + strings.Repeat("─", w-2) + "┘")
}

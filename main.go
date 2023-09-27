package main

import (
	"fmt"
	"os"

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
	logger.SetLevel(log.LevelDebug)

	logger.Info("User settings loaded...")
	logger.Infof("Log level is %s", level)

	if err := server.New(frontend.Static).Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error starting server: %s\n", err)
		os.Exit(1)
	}

	logger.Info("Exited cleanly.")
}

package main

import (
	"embed"
	"fmt"
	"github.com/ghostsecurity/reaper/backend/app"
	"github.com/ghostsecurity/reaper/backend/log"
	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/backend/workspace"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// nolint:typecheck
//
//go:embed all:frontend/dist
var assets embed.FS

func main() {

	userSettings, err := settings.Load()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading settings: %s\n", err)
		os.Exit(1)
	}

	logger := log.New(os.Stderr)
	logger.SetLevel(userSettings.Get().LogLevel)

	logger.Info("User settings loaded...")
	logger.Infof("Log level is %s", userSettings.Get().LogLevel)

	initialWorkspace, err := workspace.LoadPrevious()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading workspace: %s\n", err)
		os.Exit(1)
	}

	// Create an instance of the app structure
	application := app.New(logger.WithPrefix("app"), userSettings, initialWorkspace)

	// Create application with options
	if err := wails.Run(&options.App{
		Title:  "Reaper",
		Width:  1900,
		Height: 1024,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 0, B: 0, A: 1},
		OnStartup:        application.Startup,
		OnDomReady:       application.OnDomReady,
		OnShutdown:       application.Shutdown,
		Bind: []interface{}{
			application,
		},
		Logger: logger.WithPrefix("wails"),
	}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	logger.Info("Exited cleanly.")
}

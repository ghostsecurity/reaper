package backend

import (
	"context"
	"fmt"

	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/version"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetSettings() settings.Settings {
	return a.userSettings.Get()
}

type VersionInfo struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

func (a *App) GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: version.Version,
		Date:    version.Date,
		URL:     version.URL(),
	}
}

func (a *App) SaveSettings(newSettings *settings.Settings) {

	if err := settings.Save(newSettings); err != nil {
		a.notifyUser(fmt.Sprintf("Failed to save settings (save failed): %s", err), runtime.ErrorDialog)
		return
	}

	oldSettings := a.userSettings.Get()

	if err := a.userSettings.Modify(func(s *settings.Settings) {
		*s = *newSettings
	}); err != nil {
		a.notifyUser(fmt.Sprintf("Failed to update settings for current session: %s", err), runtime.ErrorDialog)
		return
	}

	if a.isProxyRestartRequired(oldSettings, *newSettings) {
		a.logger.Infof("Proxy settings have changed, restarting proxy...")
		if err := a.restartProxy(); err != nil {
			a.logger.Errorf("Failed to restart proxy: %s", err)
			a.notifyUser(fmt.Sprintf("Failed to restart proxy: %s", err), runtime.ErrorDialog)
			a.Shutdown(context.Background())
			return
		}
	} else {
		a.logger.Infof("Settings change does not require a proxy restart.")
	}
}

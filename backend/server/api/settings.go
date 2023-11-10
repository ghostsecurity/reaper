package api

import (
	"github.com/ghostsecurity/reaper/backend/settings"
	"github.com/ghostsecurity/reaper/version"
)

func (a *API) GetSettings() settings.Settings {
	return a.userSettings.Get()
}

type VersionInfo struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

func (a *API) GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: version.Version,
		Date:    version.Date,
		URL:     version.URL(),
	}
}

func (a *API) SaveSettings(newSettings *settings.Settings) {

	if err := settings.Save(newSettings); err != nil {
		a.notify("Failed to save settings (save failed): %s", err)
		return
	}

	oldSettings := a.userSettings.Get()

	if err := a.userSettings.Modify(func(s *settings.Settings) {
		*s = *newSettings
	}); err != nil {
		a.notify("Failed to update settings for current session: %s", err)
		return
	}

	if a.isProxyRestartRequired(oldSettings, *newSettings) {
		a.logger.Infof("Proxy settings have changed, restarting proxy...")
		if err := a.restartProxy(); err != nil {
			a.logger.Errorf("Failed to restart proxy: %s", err)
			a.notify("Failed to restart proxy: %s", err)
			return
		}
	} else {
		a.logger.Infof("Settings change does not require a proxy restart.")
	}
}

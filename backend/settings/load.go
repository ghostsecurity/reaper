package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Load() (*Provider, error) {

	configDir, err := getDir()
	if err != nil {
		return nil, fmt.Errorf("failed to access settings: %w", err)
	}

	caCert, caKey, err := LoadCA()
	if err != nil {
		return nil, fmt.Errorf("failed to load CA: %w", err)
	}

	settingsFilename := filepath.Join(configDir, settingsFilename)
	if _, err := os.Stat(settingsFilename); err != nil {
		if os.IsNotExist(err) {

			newSettings := defaultSettings.WithCA(caCert, caKey)

			// settings file doesn't exist, create it
			_ = Save(&newSettings)
			return newProvider(newSettings), nil
		}
		return nil, fmt.Errorf("failed to access settings: %w", err)
	}

	f, err := os.Open(settingsFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open settings: %w", err)
	}

	var settings Settings
	if err := json.NewDecoder(f).Decode(&settings); err != nil {
		return nil, fmt.Errorf("failed to decode settings: %w", err)
	}

	return newProvider(settings.WithCA(caCert, caKey)), nil
}

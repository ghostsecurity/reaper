package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

const (
	configDirName = ".reaper"
)

func getDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	configDir := filepath.Join(home, configDirName)
	if err := configdir.MakePath(configDir); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}
	return configDir, nil
}

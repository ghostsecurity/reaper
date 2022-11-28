package settings

import (
	"fmt"
	"github.com/kirsle/configdir"
)

const (
	configDirName = "reaper"
)

func getDir() (string, error) {
	configDir := configdir.LocalConfig(configDirName)
	if err := configdir.MakePath(configDir); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}
	return configDir, nil
}

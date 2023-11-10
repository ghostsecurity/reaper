package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Save(s *Settings) error {

	configDir, err := getDir()
	if err != nil {
		return fmt.Errorf("failed to access settings: %w", err)
	}

	settingsFilename := filepath.Join(configDir, settingsFilename)
	f, err := os.Create(settingsFilename)
	if err != nil {
		return fmt.Errorf("failed to create settings: %w", err)
	}

	if err := json.NewEncoder(f).Encode(s); err != nil {
		return fmt.Errorf("failed to encode settings: %w", err)
	}

	if saveCA([]byte(s.CACert), []byte(s.CAKey)) != nil {
		return fmt.Errorf("failed to save CA: %w", err)
	}

	return f.Close()
}

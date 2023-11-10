package settings

import (
	"github.com/ghostsecurity/reaper/backend/log"
)

const (
	settingsFilename = "config.1.json"
)

type Settings struct {
	CACert    string    `json:"ca_cert"`
	CAKey     string    `json:"ca_key"`
	ProxyPort int       `json:"proxy_port"`
	ProxyHost string    `json:"proxy_host"`
	LogLevel  log.Level `json:"log_level"`
	DarkMode  bool      `json:"dark_mode"`
}

var defaultSettings = Settings{
	ProxyPort: 8080,
	ProxyHost: "reaper",
	LogLevel:  log.LevelWarn,
	DarkMode:  true,
}

func (s Settings) WithCA(cert []byte, key []byte) Settings {
	s.CACert = string(cert)
	s.CAKey = string(key)
	return s
}

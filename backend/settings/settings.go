package settings

import "github.com/ghostsecurity/reaper/backend/log"

const (
	settingsFilename = "config.1.json"
)

type Settings struct {
	CACert    []byte    `json:"ca_cert"`
	CAKey     []byte    `json:"ca_key"`
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
	s.CACert = cert
	s.CAKey = key
	return s
}

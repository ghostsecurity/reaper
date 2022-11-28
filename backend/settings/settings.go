package settings

import "github.com/ghostsecurity/reaper/backend/log"

const (
	settingsFilename = "settings.json"
)

type Settings struct {
	CACert    []byte
	CAKey     []byte
	ProxyPort int
	ProxyHost string
	LogLevel  log.Level
	Theme     string
}

var defaultSettings = Settings{
	ProxyPort: 8080,
	ProxyHost: "reaper",
	LogLevel:  log.LevelWarn,
	Theme:     "ghost",
}

func (s Settings) WithCA(cert []byte, key []byte) Settings {
	s.CACert = cert
	s.CAKey = key
	return s
}

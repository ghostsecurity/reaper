package highlight

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

// translate reaper theme name into chroma style
func findTheme(name string) *chroma.Style {
	switch name {
	case "ghost":
		return ghostTheme
	case "dark":
		return darkTheme
	case "light":
		return lightTheme
	default:
		return styles.Fallback
	}
}

package highlight

import (
	"github.com/alecthomas/chroma/styles"
	"testing"
)

func TestThemeNamesAreLinkedToThemes(t *testing.T) {
	tests := []struct {
		name         string
		wantFallback bool
	}{
		{"ghost", false},
		{"dark", false},
		{"light", false},
		{"not a theme", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findTheme(tt.name)
			if (got == styles.Fallback) != tt.wantFallback {
				t.Errorf("findTheme() = %v, wantFallback %v", got, tt.wantFallback)
			}
			if !tt.wantFallback && got.Name != tt.name {
				t.Errorf("findTheme() = %v, want %v", got.Name, tt.name)
			}
		})
	}
}

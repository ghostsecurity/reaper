package version

import (
	"fmt"
	"strings"
)

var Version = "development"
var Date = ""

const repoURL = "https://github.com/ghostsecurity/reaper"

func URL() string {
	if strings.HasPrefix(Version, "v") {
		return fmt.Sprintf("%s/releases/tag/%s", repoURL, Version)
	}
	if Version == "development" {
		return repoURL
	}
	return fmt.Sprintf("%s/commit/%s", repoURL, Version)
}

package frontend

import "embed"

//go:embed dist/index.html
var Homepage string

//go:embed dist
var Static embed.FS

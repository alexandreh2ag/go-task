package assets

import "embed"

var (
	//go:embed templates
	TemplateFiles embed.FS
)

package tpl

import (
	"embed"
)

var (
	//go:embed login/* resource/index.tmpl
	tplFS embed.FS

	//go:embed resource
	resourceFS embed.FS

	//go:embed resource/index.tmpl
	indexPageString string
)

func GetLoginFS() embed.FS {
	return tplFS
}

func GetResourceFS() embed.FS {
	return resourceFS
}
func GetIndexPageString() string {
	return indexPageString
}

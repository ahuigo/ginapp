package config

import (
	_ "embed"
)

var (
	//go:embed conf.yaml
	ConfEmbedString string
	//go:embed conf.yaml
	ConfEmbedBytes []byte
)

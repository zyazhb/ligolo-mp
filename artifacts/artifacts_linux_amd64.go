package artifacts

import (
	"embed"
)

var (
	//go:embed agent.zip go/linux/amd64/*
	artifactsFs embed.FS
)

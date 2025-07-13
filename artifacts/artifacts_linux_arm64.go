package artifacts

import (
	"embed"
)

var (
	//go:embed agent.zip go/linux/arm64/*
	artifactsFs embed.FS
)

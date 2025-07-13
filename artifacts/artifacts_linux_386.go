package artifacts

import (
	"embed"
)

var (
	//go:embed agent.zip go/linux/386/*
	artifactsFs embed.FS
)

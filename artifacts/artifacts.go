package artifacts

import (
	"embed"
)

var (
	//go:embed agent.zip go.zip
	fs embed.FS
)

func GetGoArchive() ([]byte, error) {
	return fs.ReadFile("go.zip")
}

func GetAgentArchive() ([]byte, error) {
	return fs.ReadFile("agent.zip")
}

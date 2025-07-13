package artifacts

import (
	"path/filepath"
	"runtime"
)

func GetGoArchive() ([]byte, error) {
	return artifactsFs.ReadFile(filepath.Join("go", runtime.GOOS, runtime.GOARCH, "go.zip"))
}

func GetAgentArchive() ([]byte, error) {
	return artifactsFs.ReadFile("agent.zip")
}

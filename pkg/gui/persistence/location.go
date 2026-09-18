package persistence

import (
	"log/slog"
	"path/filepath"
)

var configFolder string

const (
	hostsFileName = "hosts.json"
)

func init() {
	err := initConfigFolder()
	if err != nil {
		slog.Error("Failed to find config location, defaulting to current directory", "error", err)
		configFolder = "."
	}
}

func HostsFile() string {
	return filepath.Join(configFolder, hostsFileName)
}

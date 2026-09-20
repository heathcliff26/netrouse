//go:build linux && !android

package persistence

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/heathcliff26/netrouse/pkg/version"
)

func initConfigFolder() {
	config, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Failed to find config location, defaulting to current directory", "error", err)
		configFolder = "./"
		return
	}
	configFolder = filepath.Join(config, strings.ToLower(version.Name))
}

//go:build linux

package persistence

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/heathcliff26/netrouse/pkg/version"
)

func initConfigFolder() error {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configFolder = xdgConfigHome
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFolder = filepath.Join(home, ".config", strings.ToLower(version.Name))
	return nil
}

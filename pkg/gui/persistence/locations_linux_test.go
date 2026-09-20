//go:build linux

package persistence

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestInitConfigFolder(t *testing.T) {
	assert := assert.New(t)

	oldFolder := ConfigFolder()
	t.Cleanup(func() {
		SetConfigFolder(oldFolder)
	})

	initConfigFolder()
	assert.Contains(configFolder, filepath.Join(".config", strings.ToLower(version.Name)), "Variable should have home location ending")

	t.Setenv("XDG_CONFIG_HOME", "/test")
	initConfigFolder()
	assert.Equal(filepath.Join("/test", strings.ToLower(version.Name)), configFolder, "Should read folder from XDG_CONFIG_HOME")
}

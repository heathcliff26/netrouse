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

	oldFolder := configFolder
	t.Cleanup(func() {
		configFolder = oldFolder
	})

	err := initConfigFolder()
	assert.NoError(err)
	assert.Contains(configFolder, filepath.Join(".config", strings.ToLower(version.Name)), "Variable should have home location ending")

	t.Setenv("XDG_CONFIG_HOME", "test")
	err = initConfigFolder()
	assert.NoError(err)
	assert.Equal("test", configFolder, "Should read folder from XDG_CONFIG_HOME")
}

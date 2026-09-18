package persistence

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostsFile(t *testing.T) {
	assert := assert.New(t)

	oldFolder := ConfigFolder()
	t.Cleanup(func() {
		SetConfigFolder(oldFolder)
	})

	configFolder = "test"
	assert.Equal(filepath.Join("test", hostsFileName), HostsFile(), "Should return the correct path")
}

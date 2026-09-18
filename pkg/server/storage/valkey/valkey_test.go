package valkey

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/heathcliff26/netrouse/pkg/server/storage/testsuite"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var newMiniredis testsuite.StorageBackendFactory = func(t *testing.T, _ string) types.StorageBackend {
	t.Helper()

	mr := miniredis.RunT(t)

	cfg := ValkeyConfig{
		Addrs: []string{mr.Addr()},
	}

	v, err := NewValkeyBackend(cfg)
	require.NoError(t, err, "Failed to create valkey backend")

	return v
}

func TestNewValkeyBackend(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		assert := assert.New(t)
		mr := miniredis.RunT(t)

		cfg := ValkeyConfig{
			Addrs: []string{mr.Addr()},
		}

		backend, err := NewValkeyBackend(cfg)
		assert.NoError(err, "Expected no error for valid config")
		assert.NotNil(backend, "Expected backend to be initialized")
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		assert := assert.New(t)
		cfg := ValkeyConfig{
			Addrs: []string{"invalid-address"},
		}

		backend, err := NewValkeyBackend(cfg)
		assert.Error(err, "Expected error for invalid config")
		assert.Nil(backend, "Expected backend to be nil for invalid config")
	})
}

func TestFileTestsuiteBasic(t *testing.T) {
	testsuite.RunStorageBackendTests(t, newMiniredis)
}

func TestFileTestsuiteRace(t *testing.T) {
	testsuite.RunStorageBackendRaceTests(t, newMiniredis)
}

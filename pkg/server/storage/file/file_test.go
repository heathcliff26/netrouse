package file

import (
	"os"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/testsuite"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

var basicTestHosts = []types.Host{
	{
		Name: "TestHost1",
		MAC:  "AA:BB:CC:DD:EE:FF",
	},
	{
		Name: "TestHost2",
		MAC:  "11:22:33:44:55:66",
	},
}

func TestNewFileBackend(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		assert := assert.New(t)
		path := "testdata/basic.yaml"

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		require.NoError(t, err, "Failed to create file backend")
		assert.NotNil(fb, "File backend should not be nil")

		assert.Equal(path, fb.path, "File backend path should match")
		assert.Equal(basicTestHosts, fb.storage.Hosts, "File backend hosts should match")
	})

	dir := t.TempDir()

	t.Run("NewHostsFile", func(t *testing.T) {
		assert := assert.New(t)
		path := dir + "/new-hosts-file.yaml"

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		require.NoError(t, err, "Failed to create file backend")
		assert.NotNil(fb, "File backend should not be nil")

		_, err = os.Stat(path)
		assert.NoError(err, "Should have new hosts file")
	})

	t.Run("UnreadableHostsFile", func(t *testing.T) {
		assert := assert.New(t)
		path := dir + "/unreadable-hosts-file.yaml"

		require.NoError(t, copyFile("testdata/basic.yaml", path, 0222), "Failed to copy file")

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		assert.Nil(fb, "File backend should be nil")
		assert.Error(err, "Should fail to create file backend with unreadable file")
		assert.Contains(err.Error(), "failed to read storage file", "Error should contain message")
	})

	t.Run("InvalidHostsFile", func(t *testing.T) {
		assert := assert.New(t)
		path := "testdata/not-yaml.txt"

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		assert.Error(err, "Should fail to create file backend with invalid YAML")
		assert.Nil(fb, "File backend should be nil")
	})

	t.Run("EnsureUppercaseMAC", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		path := dir + "/ensure-uppercase-mac-hosts-file.yaml"

		require.NoError(copyFile("testdata/lowercase.yaml", path, 0644), "Failed to copy file")

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		require.NoError(err, "Failed to create file backend")
		assert.NotNil(fb, "File backend should not be nil")

		assert.Equal(basicTestHosts, fb.storage.Hosts, "File backend hosts should match with uppercase MACs")

		f, err := os.ReadFile(path)
		require.NoError(err, "Failed to read file")
		assert.Contains(string(f), "AA:BB:CC:DD:EE:FF", "Should have written uppercase MAC to file")
	})

	t.Run("EnsureUniqueMAC", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		path := dir + "/ensure-unique-mac-hosts-file.yaml"

		require.NoError(copyFile("testdata/duplicates.yaml", path, 0644), "Failed to copy file")

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		require.NoError(err, "Failed to create file backend")
		assert.NotNil(fb, "File backend should not be nil")

		assert.Equal(basicTestHosts, fb.storage.Hosts, "File backend hosts should match with unique MACs")

		f, err := os.ReadFile(path)
		require.NoError(err, "Failed to read file")
		fs := &types.HostsFile{}
		require.NoError(yaml.Unmarshal(f, fs), "Failed to unmarshal file")
		assert.Equal(basicTestHosts, fs.Hosts, "Should have written trimmed hosts array to file")
	})

	t.Run("SaveFailure", func(t *testing.T) {
		assert := assert.New(t)
		path := dir + "/failed-to-save-hosts-file.yaml"

		require.NoError(t, copyFile("testdata/duplicates.yaml", path, 0444), "Failed to copy file")

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		assert.Error(err, "Should fail to save changed hosts file")
		assert.Nil(fb, "File backend should be nil")
		assert.Contains(err.Error(), "failed to save storage file after ensuring unique, uppercase MAC addresses:", "Error should contain message")
	})

	t.Run("OptionalAttributes", func(t *testing.T) {
		assert := assert.New(t)
		path := "testdata/optional.yaml"

		fb, err := NewFileBackend(FileBackendConfig{Path: path})
		require.NoError(t, err, "Failed to create file backend")
		assert.NotNil(fb, "File backend should not be nil")

		hosts := append([]types.Host{}, basicTestHosts...)
		hosts = append(hosts, types.Host{
			Name:    "TestHost3",
			MAC:     "77:88:99:AA:BB:CC",
			Address: "host3.example.org",
		})

		assert.Equal(path, fb.path, "File backend path should match")
		assert.Equal(hosts, fb.storage.Hosts, "File backend hosts should match")
	})
}

func TestReadonly(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	path := t.TempDir() + "/readonly-hosts-file.yaml"

	require.NoError(copyFile("testdata/basic.yaml", path, 0444), "Failed to copy file")

	fb, err := NewFileBackend(FileBackendConfig{Path: path})
	require.NoError(err, "Failed to create file backend")
	assert.NotNil(fb, "File backend should not be nil")

	readonly, err := fb.Readonly()
	assert.True(readonly, "File backend should be readonly")
	assert.NoError(err, "Should not fail to check readonly status")

	fb.path = "testdata/basic.yaml"
	readonly, err = fb.Readonly()
	assert.False(readonly, "File backend should be writable")
	assert.NoError(err, "Should not fail to check readonly status")
}

func TestFileTestsuiteBasic(t *testing.T) {
	testsuite.RunStorageBackendTests(t, newStorageBackendFactory(t))
}

func TestFileTestsuiteRace(t *testing.T) {
	testsuite.RunStorageBackendRaceTests(t, newStorageBackendFactory(t))
}

func newStorageBackendFactory(t *testing.T) testsuite.StorageBackendFactory {
	t.Helper()
	dir := t.TempDir()

	return func(t *testing.T, name string) types.StorageBackend {
		t.Helper()

		backend, err := NewFileBackend(FileBackendConfig{Path: dir + "/" + name + ".yaml"})
		require.NoError(t, err, "Failed to prepare test backend")

		return backend
	}
}

func copyFile(src, dst string, mode os.FileMode) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, mode)
}

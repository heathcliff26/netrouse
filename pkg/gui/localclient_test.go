package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/heathcliff26/netrouse/pkg/server/storage/testsuite"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalClient(t *testing.T) {
	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	t.Run("Success", func(t *testing.T) {
		require := require.New(t)
		persistence.SetConfigFolder(t.TempDir())

		c, err := newLocalClient()
		require.NoError(err, "Should succeed")
		require.NotNil(c, "Should return client")
	})
	t.Run("Mkdir", func(t *testing.T) {
		require := require.New(t)

		path := filepath.Join(t.TempDir(), "test")
		persistence.SetConfigFolder(path)

		c, err := newLocalClient()
		require.NoError(err, "Should succeed")
		require.NotNil(c, "Should return client")

		stat, err := os.Stat(path)
		require.NoError(err, "Should read stat of dir")
		require.True(stat.IsDir(), "Should create directory")
	})
	t.Run("FailMkdir", func(t *testing.T) {
		if testsuite.IsRoot() {
			t.Skip("Running as root")
		}

		require := require.New(t)

		path := filepath.Join(t.TempDir(), "foo")
		err := os.MkdirAll(path, 0444)
		require.NoError(err, "Should create dir as readonly")
		path = filepath.Join(path, "bar")
		persistence.SetConfigFolder(path)

		c, err := newLocalClient()
		require.Error(err, "Should fail")
		require.NotNil(c, "Should return client")
	})
	t.Run("FailReadHosts", func(t *testing.T) {
		require := require.New(t)

		path := "testdata/not-yaml.txt"
		persistence.SetConfigFolder(path)

		c, err := newLocalClient()
		require.Error(err, "Should fail")
		require.NotNil(c, "Should return client")
	})
}

func TestStorageFunctions(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})
	persistence.SetConfigFolder(t.TempDir())

	c, err := newLocalClient()
	require.NoError(err, "Should succeed")
	require.NotNil(c, "Should return client")

	host := types.Host{
		Name: "TestHost1",
		MAC:  "AA:BB:CC:DD:EE:FF",
	}

	err = c.AddHost(t.Context(), host)
	require.NoError(err, "Should add host")

	hosts, err := c.GetHosts(t.Context())
	require.NoError(err, "Should get hosts")
	require.Equal(1, len(hosts), "Should have one host")
	require.Equal(host, hosts[0], "Should have correct host")

	err = c.RemoveHost(t.Context(), host.MAC)
	require.NoError(err, "Should remove host")

	hosts, err = c.GetHosts(t.Context())
	require.NoError(err, "Should get hosts")
	require.Equal(0, len(hosts), "Should have no hosts")
}

func TestStatus(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})
	persistence.SetConfigFolder(t.TempDir())
	c, err := newLocalClient()
	require.NoError(err, "Should succeed")
	require.NotNil(c, "Should return client")

	h1 := types.Host{
		Name:    "TestHost1",
		MAC:     "00:11:22:33:44:55",
		Address: "127.0.0.1",
	}
	err = c.AddHost(t.Context(), h1)
	require.NoError(err, "Should add host")
	h2 := types.Host{
		Name: "TestHost2",
		MAC:  "AA:BB:CC:DD:EE:FF",
	}
	err = c.AddHost(t.Context(), h2)
	require.NoError(err, "Should add host")

	status, err := c.Status(t.Context())
	require.NoError(err, "Should get status")
	require.Len(status, 1, "Should have one statuses")
}

func TestWake(t *testing.T) {
	tMatrix := []struct {
		Name, MAC string
		Error     string
	}{
		{
			Name: "Success",
			MAC:  "ff:ff:ff:ff:ff:ff",
		},
		{
			Name:  "InvalidMAC",
			MAC:   "Not-a-mac-address",
			Error: "failed to parse MAC address",
		},
	}

	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})
	persistence.SetConfigFolder(t.TempDir())
	c, err := newLocalClient()
	require.NoError(err, "Should succeed")
	require.NotNil(c, "Should return client")

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			err := c.Wake(t.Context(), tCase.MAC)
			if tCase.Error != "" {
				assert.ErrorContains(err, tCase.Error)
			} else {
				assert.NoError(err)
			}
		})
	}
}

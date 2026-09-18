package testsuite

import (
	"fmt"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Runs basic tests for the storage backend.
// It will create a new backend instance for each test case.
func RunStorageBackendTests(t *testing.T, factory StorageBackendFactory) {
	t.Run("AddHost", func(t *testing.T) {
		backend := factory(t, "add-host")

		err := backend.AddHost(testHosts[0])
		require.NoError(t, err, "AddHost failed")
	})

	t.Run("GetHost", func(t *testing.T) {
		backend := factory(t, "get-host")
		addHosts(t, backend)

		host, err := backend.GetHost(testHosts[0].MAC)
		require.NoError(t, err, "GetHost failed")
		assert.Equal(t, testHosts[0], host, "Failed to retrieve host")
	})

	t.Run("RemoveHost", func(t *testing.T) {
		tMatrix := []struct {
			name  string
			index int
		}{
			{
				name:  "FirstElement",
				index: 0,
			},
			{
				name:  "LastElement",
				index: len(testHosts) - 1,
			},
			{
				name:  "MiddleElement",
				index: 1,
			},
		}

		for i, tCase := range tMatrix {
			t.Run(tCase.name, func(t *testing.T) {
				backend := factory(t, fmt.Sprintf("remove-host-%d", i))
				addHosts(t, backend)

				err := backend.RemoveHost(testHosts[tCase.index].MAC)
				require.NoError(t, err, "RemoveHost failed")
				host, _ := backend.GetHost(testHosts[tCase.index].MAC)
				assert.Empty(t, host, "Expected empty host")
			})
		}
	})

	t.Run("RemoveHostNonExistent", func(t *testing.T) {
		backend := factory(t, "remove-host-non-existent")

		err := backend.RemoveHost("00:11:22:33:44:55")
		require.NoError(t, err, "RemoveHost for non-existent host failed")
	})

	t.Run("GetHosts", func(t *testing.T) {
		backend := factory(t, "get-hosts")
		addHosts(t, backend)

		hosts, err := backend.GetHosts()
		require.NoError(t, err, "GetHosts failed")
		assert.Equal(t, testHosts, hosts, "Expected result to match testHosts")
	})

	t.Run("CaseInsensitiveMAC", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		backend := factory(t, "case-insensitive-mac")

		err := backend.AddHost(types.Host{MAC: "aa:bb:cc:dd:ee:ff", Name: "LowerCase"})

		require.NoError(err, "Should add host")

		host, err := backend.GetHost("AA:BB:CC:dd:ee:ff")
		require.NoError(err, "Should get host regardless of case")
		assert.Equal("LowerCase", host.Name, "Should get correct host name")

		hosts, err := backend.GetHosts()
		require.NoError(err, "Should get hosts")

		require.Len(hosts, 1, "Should have one host")
		assert.Equal("LowerCase", hosts[0].Name, "Should get correct host name")
		assert.Equal("AA:BB:CC:DD:EE:FF", hosts[0].MAC, "Should store MAC in upper case")
	})

	t.Run("AddHostOverwrite", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		backend := factory(t, "add-host-overwrite")

		err := backend.AddHost(testHosts[0])
		require.NoError(err, "Should add host")

		host := testHosts[0]
		host.Name = "NewName"
		err = backend.AddHost(host)
		require.NoError(err, "Should overwrite host")

		hosts, err := backend.GetHosts()
		require.NoError(err, "Should get hosts")
		require.Len(hosts, 1, "Should have one host")
		assert.Equal("NewName", hosts[0].Name, "Should have new name")
	})

	t.Run("HostOverwriteShouldKeepOrder", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)
		backend := factory(t, "host-overwrite-keep-order")
		addHosts(t, backend)

		testHost := testHosts[0]
		testHost.Name = "NewName"
		err := backend.AddHost(testHost)
		require.NoError(err, "Should overwrite host")

		hosts, err := backend.GetHosts()
		require.NoError(err, "Should get hosts")
		require.Len(hosts, len(testHosts), "Should have same number of hosts")
		assert.Equal("NewName", hosts[0].Name, "Should have new name")

		for i, host := range hosts {
			assert.Equal(testHosts[i].MAC, host.MAC, "Should keep order")
		}
	})
}

package testsuite

import (
	"os/user"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/require"
)

type StorageBackendFactory func(t *testing.T, name string) types.StorageBackend

var testHosts = []types.Host{
	{
		MAC:  "AA:BB:CC:DD:EE:FF",
		Name: "TestHost1",
	},
	{
		MAC:     "11:22:33:44:55:66",
		Name:    "TestHost2",
		Address: "host.example.com",
	},
	{
		MAC:  "77:88:99:AA:BB:CC",
		Name: "TestHost3",
	},
	{
		MAC:  "FF:88:99:AA:BB:CC",
		Name: "TestHost4",
	},
	{
		MAC:  "FE:11:99:AA:BB:CC",
		Name: "TestHost5",
	},
}

func addHosts(t *testing.T, backend types.StorageBackend) {
	t.Helper()

	for _, host := range testHosts {
		err := backend.AddHost(host)
		require.NoError(t, err, "AddHost failed for %s", host.Name)
	}
}

func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

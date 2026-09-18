package valkey

import (
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeHostBackwardsCompatibility(t *testing.T) {
	assert := assert.New(t)

	host := types.Host{
		MAC:  "AA:BB:CC:DD:EE:FF",
		Name: "TestHost",
	}

	result := deserializeHost(host.MAC, host.Name)
	assert.Equal(host, result, "Deserialized host should match original host")
}

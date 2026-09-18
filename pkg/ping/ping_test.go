package ping

import (
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
)

func TestUnresolvableAddress(t *testing.T) {
	assert := assert.New(t)

	hosts := []types.Host{
		{
			MAC:     "00:11:22:33:44:55",
			Address: "unresolvable.domain",
		},
	}

	result := PingHosts(hosts)

	assert.Equal(hosts[0].MAC, result[0].MAC, "MAC address should match")
	assert.Equal(hosts[0].Address, result[0].Address, "Address should match")
	assert.False(result[0].Online, "Host should be offline for unresolvable address")
	assert.Contains(result[0].Error, "no such host", "Error should indicate unresolvable address")
}

func TestPing(t *testing.T) {
	t.Run("IPv4", func(t *testing.T) {
		assert := assert.New(t)

		hosts := []types.Host{
			{
				MAC:     "00:11:22:33:44:55",
				Address: "127.0.0.1",
			},
		}
		result := PingHosts(hosts)

		assert.True(result[0].Online, "Host should be online")
		assert.Empty(result[0].Error, "Should not return an error")
	})
	t.Run("IPv6", func(t *testing.T) {
		assert := assert.New(t)

		hosts := []types.Host{
			{
				MAC:     "00:11:22:33:44:55",
				Address: "::1",
			},
		}
		result := PingHosts(hosts)

		assert.True(result[0].Online, "Host should be online")
		assert.Empty(result[0].Error, "Should not return an error")
	})
	t.Run("UnreachableHost", func(t *testing.T) {
		assert := assert.New(t)

		hosts := []types.Host{
			{
				MAC:     "00:11:22:33:44:55",
				Address: "192.0.2.1", // TEST-NET-1 IP address, should be unreachable
			},
		}
		result := PingHosts(hosts)

		assert.False(result[0].Online, "Host should be offline")
		assert.Empty(result[0].Error, "Should not return an error")
	})
}

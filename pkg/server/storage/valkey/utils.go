package valkey

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
)

const (
	keyName    = "name"
	keyAddress = "address"
)

// serialize the Host so it can be stored as a single value in Valkey.
func serializeHost(host types.Host) string {
	result := fmt.Sprintf("%s=%s;", keyName, host.Name)
	if host.Address != "" {
		result += fmt.Sprintf("%s=%s;", keyAddress, host.Address)
	}
	return result
}

// Deserialize host back into it's parameters.
func deserializeHost(mac, data string) types.Host {
	keyvalues := strings.Split(data, ";")

	host := types.Host{
		MAC: mac,
	}

	// Maintain compatibility with older versions that only stored the name directly.
	if len(keyvalues) == 1 {
		host.Name = keyvalues[0]
		return host
	}

	for i, kv := range keyvalues {
		if i == len(keyvalues)-1 && kv == "" {
			// Skip the last empty element if it exists
			continue
		}

		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			slog.Warn("Received invalid host data from valkey, expected pairs key=value separated by semicolons", slog.String("data", data))
			continue
		}
		switch pair[0] {
		case keyName:
			host.Name = pair[1]
		case keyAddress:
			host.Address = pair[1]
		default:
			slog.Warn("Received unknown key in host data from valkey", slog.String("key", pair[0]), slog.String("data", data))
		}
	}
	return host
}

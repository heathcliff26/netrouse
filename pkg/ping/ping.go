package ping

import (
	"time"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	probing "github.com/prometheus-community/pro-bing"
)

// Ping all the hosts and return their status.
// Will return the first error encountered, if any.
func PingHosts(hosts []types.Host) []types.HostStatus {
	res := make(chan types.HostStatus, 1)

	for _, host := range hosts {
		status := types.HostStatus{
			MAC:     host.MAC,
			Address: host.Address,
		}

		go func() {
			pinger, err := probing.NewPinger(host.Address)
			if err != nil {
				status.Error = err.Error()
				res <- status
				return
			}

			pinger.Count = 1
			pinger.Timeout = 200 * time.Millisecond

			err = pinger.Run()
			if err != nil {
				status.Error = err.Error()
				res <- status
				return
			}
			stats := pinger.Statistics()

			status.Online = stats != nil && stats.PacketsRecv > 0
			res <- status
		}()
	}

	results := make([]types.HostStatus, len(hosts))
	for i := range results {
		status := <-res
		results[i] = status
	}
	return results
}

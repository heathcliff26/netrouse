//go:build race

package testsuite

import (
	"fmt"
	"sync"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
)

// These tests do not run any checks themselves, they rely on race detection for that.
func RunStorageBackendRaceTests(t *testing.T, factory StorageBackendFactory) {
	t.Run("ConcurrentAddHost", func(t *testing.T) {
		var wg sync.WaitGroup
		backend := factory(t, "concurrent-add-host")

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				mac := fmt.Sprintf("AA:BB:CC:DD:EE:%02X", i)
				host := fmt.Sprintf("TestHost%d", i)
				_ = backend.AddHost(types.Host{
					MAC:  mac,
					Name: host,
				})
			}(i)
		}
		wg.Wait()
	})

	t.Run("ConcurrentGetHost", func(t *testing.T) {
		var wg sync.WaitGroup
		backend := factory(t, "concurrent-get-host")
		addHosts(t, backend)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, _ = backend.GetHost(testHosts[0].MAC)
			}()
		}
		wg.Wait()
	})

	t.Run("ConcurrentRemoveHost", func(t *testing.T) {
		var wg sync.WaitGroup
		backend := factory(t, "concurrent-remove-host")
		addHosts(t, backend)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_ = backend.RemoveHost(testHosts[0].MAC)
			}()
		}
		wg.Wait()
	})

	t.Run("ConcurrentGetHosts", func(t *testing.T) {
		var wg sync.WaitGroup
		backend := factory(t, "concurrent-get-hosts")
		addHosts(t, backend)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, _ = backend.GetHosts()
			}()
		}
		wg.Wait()
	})
}

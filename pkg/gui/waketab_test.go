package gui

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWakeTab(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")

	tab := newLocalTab(w)
	require.NotNil(tab, "Tab should not be nil")
	assert.Equal(w, tab.window, "Should have parent window")

	assert.NotNil(tab.client, "Tab should have client")
	assert.NotNil(tab.tab.Content, "Tab should have content")

	assert.NotNil(tab.title, "Should have title")
	assert.NotNil(tab.errFetch, "Should have errFetch")
	assert.NotNil(tab.errStatus, "Should have errStatus")
	assert.NotNil(tab.hostsContainer, "Should have hosts container")
}

func TestWakeTabUpdate(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")

	tab := newLocalTab(w)
	require.NotNil(tab)
	require.NotNil(tab.client)
	require.NotNil(tab.tab.Content)

	require.Empty(tab.hostsContainer.Objects, "Hosts should be empty")

	host := types.Host{
		Name: "Test Host",
		MAC:  "00:11:22:33:44:55",
	}
	require.NoError(tab.client.AddHost(host), "Should add host")

	tab.fetchHosts()

	require.Len(tab.hostsContainer.Objects, 1, "Should have added host")
	require.Equal(tab.hosts[0].card, tab.hostsContainer.Objects[0], "Host objects should match")
}

func TestFetchHost(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")

	tab := newLocalTab(w)
	require.NotNil(tab)
	require.NotNil(tab.client)
	require.NotNil(tab.tab.Content)

	host := types.Host{
		Name:    "Test Host",
		MAC:     "00:11:22:33:44:55",
		Address: "127.0.0.1",
	}
	require.NoError(tab.client.AddHost(host), "Should add host")

	tab.errFetch.Show()
	tab.errStatus.Show()
	tab.fetchHosts()
	require.Len(tab.hosts, 1, "Should have fetched hosts")
	require.True(tab.errFetch.Hidden, "Should hide fetch error")

	acquired := false
	deadline := time.Now().Add(1 * time.Second)

	time.Sleep(5 * time.Millisecond)
	for time.Now().Before(deadline) {
		if tab.lock.TryLock() {
			acquired = true
			t.Cleanup(tab.lock.Unlock)
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	require.True(acquired, "Should have acquired lock")
	require.True(tab.errStatus.Hidden, "Should hide status error")
}

func TestWakeTabSelected(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")

	tab := newLocalTab(w)
	require.NotNil(tab)
	require.NotNil(tab.client)
	require.NotNil(tab.tab.Content)

	require.Error(tab.ctx.Err(), "Context should start with error")
	tab.selected()
	require.NoError(tab.ctx.Err(), "Context should not have error")

	ctx := tab.ctx
	tab.selected()
	require.Equal(ctx, tab.ctx, "Consecutive calls to selected should be noop")

	tab.unselected()
	require.Error(tab.ctx.Err(), "Context should have error")
}

func TestWakeTabSetRemote(t *testing.T) {
	require := require.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")

	tab := newTabFromRemote(w, &persistence.RemoteServer{
		Name: "Test",
		URL:  "localhost",
	})
	require.NotNil(tab, "Should have created tab")
	require.NotNil(tab.client, "Should have client")
	require.NotNil(tab.remote, "Should have remote")

	tab.client = nil
	remote := persistence.RemoteServer{
		Name: "Changed",
		URL:  "not-a-host.local",
	}
	tab.SetRemote(remote)
	require.Equal(remote, *tab.remote, "Should have updated remote")
	require.Contains(tab.title.Text, remote.Name, "Should have updated title")
	require.NotNil(tab.client, "Should have updated client")
}

func TestHostWidget(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	oldFolder := persistence.ConfigFolder()
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
	})

	persistence.SetConfigFolder(t.TempDir())
	app := test.NewApp()
	w := app.NewWindow("Test")
	tab := newLocalTab(w)

	host := types.Host{
		Name:    "Test Host",
		MAC:     "00:11:22:33:44:55",
		Address: "127.0.0.1",
	}
	err := tab.client.AddHost(host)
	require.NoError(err, "Should add host")
	tab.fetchHosts()
	require.Len(tab.hosts, 1, "Should have added host")

	// Lock tab to prevent race with tab.updateStatus()
	tab.lock.Lock()

	hostWidget := tab.hosts[0]
	require.NotNil(hostWidget, "Should have created host widget")
	require.Equal(hostWidget.host, host, "Should have set host")
	require.NotNil(hostWidget.card, "Should have created card")
	require.NotNil(hostWidget.status, "Should have created status")
	require.NotNil(hostWidget.address, "Should have created address")
	require.NotNil(hostWidget.deleteBtn, "Should have created delete button")
	require.NotNil(hostWidget.wakeBtn, "Should have created wake button")
	require.NotNil(hostWidget.statusContainer, "Should have created status container")
	require.False(hostWidget.statusContainer.Hidden, "Should have visible status")

	assert.Equal(hostStatusUnknown, hostWidget.status.Resource, "Should have unkown status at start")

	hostWidget.updateStatus(types.HostStatus{
		MAC:    host.MAC,
		Online: false,
	})
	assert.Equal(hostStatusOffline, hostWidget.status.Resource, "Should have offline status")

	hostWidget.updateStatus(types.HostStatus{
		MAC:    host.MAC,
		Online: false,
		Error:  "Test",
	})
	assert.Equal(hostStatusUnknown, hostWidget.status.Resource, "Should have unkown status")

	hostWidget.updateStatus(types.HostStatus{
		MAC:    host.MAC,
		Online: true,
	})
	assert.Equal(hostStatusOnline, hostWidget.status.Resource, "Should have online status")

	// Unlock tab to allow editing host
	tab.lock.Unlock()

	host.Name = "Changed"
	host.Address = ""
	err = tab.client.AddHost(host)
	require.NoError(err, "Should edit host")
	tab.fetchHosts()
	require.Len(tab.hosts, 1, "Should have added host")
	require.Same(hostWidget, tab.hosts[0], "Should have updated widget in place")

	// Lock tab to prevent race with tab.updateStatus()
	tab.lock.Lock()

	assert.Equal(host, hostWidget.host, "Should have updated host")
	assert.Equal(host.Name, hostWidget.card.Title, "Should have updated title")
	assert.True(hostWidget.statusContainer.Hidden, "Should have hidden status")
	assert.Equal(host.Address, hostWidget.address.Text, "Should have updated address")
	assert.Equal(hostStatusUnknown, hostWidget.status.Resource, "Should have changed status to unknown")

	// Unlock tab to allow removing host
	tab.lock.Unlock()

	hostWidget.deleteBtn.OnTapped()

	assert.Len(tab.hosts, 0, "Should have removed host")
}

package gui

import (
	"testing"

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
	require.Equal(tab.hosts[0].object, tab.hostsContainer.Objects[0], "Host objects should match")
}

// TODO: Test fetch and status

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

	hostWidget := tab.hosts[0]
	require.NotNil(hostWidget)
	require.NotNil(hostWidget.object)
	require.NotNil(hostWidget.status)

	hostWidget.updateStatus(types.HostStatus{MAC: host.MAC, Online: true})
	assert.NotNil(hostWidget.status)

	hostWidget.deleteBtn.OnTapped()

	require.Len(tab.hosts, 0, "Should have removed host")
}

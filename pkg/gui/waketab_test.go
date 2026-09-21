package gui

import (
	"context"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/lang"
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
	require.NotNil(tab)
	require.NotNil(tab.client)
	assert.NotNil(tab.tab.Content)

	tab.selected()
	require.NotNil(tab.ctx)
	assert.NotNil(tab.cancel)

	tab.unselected()
	assert.Equal(context.Canceled, tab.ctx.Err())

	host := types.Host{
		Name:    "Test Host",
		MAC:     "00:11:22:33:44:55",
		Address: "127.0.0.1",
	}
	hostWidget := newHostWidget(tab, host)
	require.NotNil(hostWidget)
	require.NotNil(hostWidget.object)
	require.NotNil(hostWidget.status)

	hostWidget.updateStatus(types.HostStatus{MAC: host.MAC, Online: true})
	assert.NotNil(hostWidget.status)
}

func TestWakeTabUpdate(t *testing.T) {
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
	require.NotNil(tab)
	require.NotNil(tab.client)
	require.NotNil(tab.tab.Content)

	hosts := extractHostObjectsFromTab(t, tab)
	require.Empty(hosts, "Hosts should be empty")

	host := types.Host{
		Name: "Test Host",
		MAC:  "00:11:22:33:44:55",
	}
	require.NoError(tab.client.AddHost(host), "Should add host")

	tab.fetchHosts()

	hosts = extractHostObjectsFromTab(t, tab)
	require.Len(hosts, 1)
	assert.Equal(tab.hosts[0].object, hosts[0])

	tab.errFetch = true
	tab.errStatus = true
	tab.update()

	hosts = extractHostObjectsFromTab(t, tab)
	require.Len(hosts, 3, "Should have added errors to hosts")
	assert.Equal(tab.hosts[0].object, hosts[0])

	errorText, ok := hosts[1].(*canvas.Text)
	require.True(ok, "Second object should be an error text")
	assert.Contains(errorText.Text, lang.L("error.fetchHosts"))

	errorText, ok = hosts[2].(*canvas.Text)
	require.True(ok, "Third object should be an error text")
	assert.Contains(errorText.Text, lang.L("error.getStatus"))

	tab.fetchHosts()
	hosts = extractHostObjectsFromTab(t, tab)
	assert.Len(hosts, 2, "FetchHosts should have removed fetch error")

	tab.updateStatus()
	hosts = extractHostObjectsFromTab(t, tab)
	assert.Len(hosts, 1, "UpdateStatus should have removed status error")
}

func extractHostObjectsFromTab(t *testing.T, tab *wakeTab) []fyne.CanvasObject {
	t.Helper()
	require := require.New(t)

	border, ok := tab.tab.Content.(*fyne.Container)
	require.True(ok, "Tab content should be a container")
	require.Len(border.Objects, 3, "Border should have center, top, and bottom objects")
	hosts, ok := border.Objects[0].(*fyne.Container)
	require.True(ok, "Hosts should be a container")

	return hosts.Objects
}

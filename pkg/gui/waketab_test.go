package gui

import (
	"context"
	"testing"

	fApp "fyne.io/fyne/v2/app"
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
	newApp = test.NewApp
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
		newApp = fApp.New
	})

	persistence.SetConfigFolder(t.TempDir())
	app := New()

	tab := app.tabLocal
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

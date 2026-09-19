package gui

import (
	"testing"

	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoteCard(t *testing.T) {
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

	card := NewRemoteCard(app)
	require.NotNil(card)
	require.NotNil(card.CreateRenderer())

	remote := persistence.RemoteServer{Name: "Primary", URL: "http://example.com"}
	card.SetRemote(remote)
	assert.Equal(remote.Name, card.card.Title)
	assert.Equal(remote.URL, card.card.Subtitle)

	app.settings.Remotes = []persistence.RemoteServer{remote}
	app.addRemote(&remote)
	require.Len(app.tabs, 1)
	assert.Len(app.settings.Remotes, 1)

	card.remote = remote
	card.delete()

	assert.Len(app.settings.Remotes, 0)
	assert.Len(app.tabs, 0)
	assert.Len(app.appTabs.Items, 2)
}

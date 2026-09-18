package gui

import (
	"testing"

	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/heathcliff26/netrouse/pkg/gui/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	oldFolder := persistence.ConfigFolder()
	newApp = test.NewApp
	t.Cleanup(func() {
		persistence.SetConfigFolder(oldFolder)
		newApp = fApp.New
	})

	app := New()

	require.NotNil(t, app)
	assert.NotNil(app.app)
	assert.NotNil(app.main)
	assert.NotNil(app.tabLocal)
	assert.NotNil(app.tabSettings)
	assert.Nil(app.tabs)
}

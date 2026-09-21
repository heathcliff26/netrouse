package customthemes

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusIconTheme(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(statusIcon{}, NewStatusIconTheme())
}

func TestStatusIcon(t *testing.T) {
	assert := assert.New(t)

	// Themes need an app to run
	_ = test.NewApp()

	assert.Equal(colorGreen, statusIcon{}.Color(theme.ColorNamePrimary, theme.VariantDark))
	assert.Equal(colorRed, statusIcon{}.Color(theme.ColorNameError, theme.VariantDark))
	assert.Equal(colorGray, statusIcon{}.Color(theme.ColorNameDisabled, theme.VariantDark))

	assert.Equal(theme.DefaultTheme().Font(fyne.TextStyle{}), statusIcon{}.Font(fyne.TextStyle{}))
	assert.Equal(theme.DefaultTheme().Icon(theme.IconNameAccount), statusIcon{}.Icon(theme.IconNameAccount))
	assert.Equal(theme.DefaultTheme().Size(theme.SizeNameSeparatorThickness), statusIcon{}.Size(theme.SizeNameSeparatorThickness))
	assert.Equal(theme.DefaultTheme().Color(theme.ColorNameBackground, theme.VariantDark), statusIcon{}.Color(theme.ColorNameBackground, theme.VariantDark))
}

func TestColorRed(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(colorRed, Red())
}

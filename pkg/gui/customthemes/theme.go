package customthemes

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = statusIcon{}

var (
	colorGreen = color.RGBA{R: 0, G: 160, B: 0, A: math.MaxUint8}
	colorRed   = color.RGBA{R: 200, G: 0, B: 0, A: math.MaxUint8}
	colorGray  = color.Gray{100}
)

type statusIcon struct{}

func NewStatusIconTheme() fyne.Theme {
	return statusIcon{}
}

func (statusIcon) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (statusIcon) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (statusIcon) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (statusIcon) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return colorGreen
	case theme.ColorNameError:
		return colorRed
	case theme.ColorNameDisabled:
		return colorGray
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

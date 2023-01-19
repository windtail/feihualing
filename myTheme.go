package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type myTheme struct {
}

//go:embed NotoSansSC-Regular.ttf
var notoScContent []byte

//go:embed NotoSansSC-Bold.ttf
var notoScBoldContent []byte

var notoSc = &fyne.StaticResource{
	StaticName:    "NotoSansSC-Regular.ttf",
	StaticContent: notoScContent,
}

var notoScBold = &fyne.StaticResource{
	StaticName:    "NotoSansSC-Bold.ttf",
	StaticContent: notoScBoldContent,
}

var _ fyne.Theme = (*myTheme)(nil)

func (t *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}
func (t *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return notoScBold
	}

	return notoSc
}
func (t *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t *myTheme) Size(size fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(size)
}

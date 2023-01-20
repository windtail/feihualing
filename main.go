package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("飞花令")
	myApp.Settings().SetTheme(&myTheme{})
	mgr := NewScreenManager()

	if _, ok := myApp.(desktop.App); ok {
		myWindow.Resize(fyne.NewSize(800, 600))
	}

	poems := NewPoems()
	poems.LoadDefault()

	mgr.Add("detail", NewDetailScreen(poems, mgr, myWindow))
	mgr.Add("entry", NewEntryScreen(poems, mgr, myWindow))

	myWindow.SetContent(mgr.Build("entry"))
	myWindow.ShowAndRun()
}

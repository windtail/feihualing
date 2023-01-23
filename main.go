package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
)

func loadPoems(dir fyne.URI) (*Poems, error) {
	poems := NewPoems()
	if path, err := storage.Child(dir, "poems.db"); err != nil {
		return nil, err
	} else {
		if err := poems.Init(path.String()); err != nil {
			return nil, err
		} else {
			return poems, nil
		}
	}
}

func main() {
	myApp := app.NewWithID("cn.poem.flower")
	myWindow := myApp.NewWindow("飞花令")
	myApp.Settings().SetTheme(&myTheme{})
	mgr := NewScreenManager()

	if _, ok := myApp.(desktop.App); ok {
		myWindow.Resize(fyne.NewSize(800, 600))
	}

	poems, err := loadPoems(myApp.Storage().RootURI())
	if err != nil {
		panic(err.Error())
	}

	mgr.Add("detail", NewDetailScreen(poems, mgr, myWindow))
	mgr.Add("entry", NewEntryScreen(poems, mgr, myWindow))

	myWindow.SetContent(mgr.Build("entry"))
	myWindow.ShowAndRun()
}

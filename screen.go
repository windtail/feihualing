package main

import "fyne.io/fyne/v2"

type ScreenManager struct {
	Screens map[string]fyne.CanvasObject
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{Screens: make(map[string]fyne.CanvasObject)}
}

func (mgr *ScreenManager) Add(name string, screen fyne.CanvasObject) {
	mgr.Screens[name] = screen
}

func (mgr *ScreenManager) SwitchTo(name string) {
	if s, ok := mgr.Screens[name]; ok {
		for _, screen := range mgr.Screens {
			if s == screen {
				screen.Show()
			} else {
				screen.Hide()
			}
		}
	}
}

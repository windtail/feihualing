package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Screen interface {
	Show(interface{})
	Hide()
	RootObj() fyne.CanvasObject
}

type ScreenManager struct {
	Screens map[string]Screen
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{Screens: make(map[string]Screen)}
}

func (mgr *ScreenManager) Add(name string, screen Screen) {
	mgr.Screens[name] = screen
}

func (mgr *ScreenManager) SwitchTo(name string) {
	mgr.SwitchToWithCtx(name, nil)
}

func (mgr *ScreenManager) SwitchToWithCtx(name string, ctx interface{}) {
	if s, ok := mgr.Screens[name]; ok {
		for _, screen := range mgr.Screens {
			if s == screen {
				screen.Show(ctx)
			} else {
				screen.Hide()
			}
		}
	}
}

func (mgr *ScreenManager) Build(entry string) fyne.CanvasObject {
	mgr.SwitchTo(entry)

	objs := make([]fyne.CanvasObject, 0, len(mgr.Screens))
	for _, s := range mgr.Screens {
		objs = append(objs, s.RootObj())
	}

	return container.NewMax(objs...)
}

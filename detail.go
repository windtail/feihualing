package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DetailContext struct {
	poem   *Poem
	search *Search
}

func NewDetailContext(poem *Poem, search *Search) *DetailContext {
	return &DetailContext{poem: poem, search: search}
}

type DetailScreen struct {
	root fyne.CanvasObject
	ctx  binding.Untyped
}

func NewDetailScreen(poems *Poems, mgr *ScreenManager, win fyne.Window) *DetailScreen {
	context := binding.NewUntyped()

	text := widget.NewRichTextWithText("")

	returnBtn := widget.NewButtonWithIcon("返回", theme.NavigateBackIcon(), func() {
		mgr.SwitchTo("entry")
	})
	delBtn := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {
		if ctx, err := context.Get(); err != nil || ctx == nil {
			return
		} else {
			p := ctx.(*DetailContext).poem
			dialog.ShowConfirm("警告", fmt.Sprintf("删除 %s ？", p.Title), func(b bool) {
				if !b {
					return
				}

				if err := poems.Remove(p); err != nil {
					dialog.ShowError(err, win)
				} else {
					mgr.SwitchToWithCtx("entry", true)
				}
			}, win)
		}
	})
	editBtn := widget.NewButtonWithIcon("编辑", theme.DocumentCreateIcon(), func() {
		if ctx, err := context.Get(); err != nil || ctx == nil {
			return
		} else {
			p := ctx.(*DetailContext).poem
			mgr.SwitchToWithCtx("edit", NewEditContext(p))
		}
	})

	context.AddListener(binding.NewDataListener(func() {
		ctx, err := context.Get()
		if err != nil || ctx == nil {
			return
		}
		p := ctx.(*DetailContext)

		text.ParseMarkdown(p.poem.DetailMarkdown(p.search))
	}))

	return &DetailScreen{
		root: container.NewBorder(nil, container.NewGridWithColumns(3, returnBtn, editBtn, delBtn), nil, nil, container.NewScroll(text)),
		ctx:  context,
	}
}

func (s *DetailScreen) Show(ctx interface{}) {
	c := ctx.(*DetailContext)
	if c.search == nil { // 保持search不变
		if old, err := s.ctx.Get(); err != nil || old == nil {
			c.search = EmptySearch()
		} else {
			c.search = old.(*DetailContext).search
		}
	}

	_ = s.ctx.Set(ctx)
	s.root.Show()
}

func (s *DetailScreen) Hide() {
	s.root.Hide()
}

func (s *DetailScreen) RootObj() fyne.CanvasObject {
	return s.root
}

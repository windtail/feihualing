package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

type EditContext struct {
	poem *Poem
}

func NewEditContext(poem *Poem) *EditContext {
	return &EditContext{poem: poem}
}

type EditScreen struct {
	root fyne.CanvasObject
	ctx  binding.Untyped
}

func NewEditScreen(poems *Poems, mgr *ScreenManager, win fyne.Window) *EditScreen {
	context := binding.NewUntyped()

	no := widget.NewEntry()
	no.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return nil
		}

		if _, err := strconv.ParseUint(s, 10, 64); err != nil {
			return errors.New("请输入正整数")
		}
		return nil
	}
	title := widget.NewEntry()
	title.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return errors.New("标题不能为空白")
		}
		return nil
	}
	author := widget.NewEntry()
	author.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return errors.New("作者不能为空白")
		}
		return nil
	}
	dynasty := widget.NewEntry()
	dynasty.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return errors.New("朝代不能为空白")
		}
		return nil
	}
	content := widget.NewMultiLineEntry()
	content.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return errors.New("内容不能为空白")
		}
		return nil
	}

	editable := container.NewBorder(container.New(layout.NewFormLayout(), widget.NewLabel("序号"), no,
		widget.NewLabel("标题"), title,
		widget.NewLabel("朝代"), dynasty,
		widget.NewLabel("作者"), author), nil, nil, nil, content)

	saveBtn := widget.NewButtonWithIcon("保存", theme.DocumentSaveIcon(), func() {
		if ctx, err := context.Get(); err != nil || ctx == nil {
			return
		} else {
			err = no.Validate()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			err = title.Validate()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			err = author.Validate()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			err = dynasty.Validate()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			err = content.Validate()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			no_, title_, author_, dynasty_, content_ := strings.TrimSpace(no.Text), strings.TrimSpace(title.Text), strings.TrimSpace(author.Text), strings.TrimSpace(dynasty.Text), strings.TrimSpace(content.Text)

			p := ctx.(*EditContext).poem
			if p != nil { // switch from detail, switch back
				if len(no_) == 0 {
					dialog.ShowError(errors.New("编辑时序号不能为空"), win)
					return
				}
				aNo, _ := strconv.ParseUint(no_, 10, 64)
				aPoem := NewPoem(aNo, title_, dynasty_, author_, content_)
				if err := poems.Modify(p, aPoem); err != nil {
					dialog.ShowError(err, win)
					return
				}
				mgr.SwitchToWithCtx("detail", NewDetailContext(aPoem, nil))
			} else {
				var aNo uint64 = 0
				if len(no_) == 0 {
					aNo = poems.NextNo()
				} else {
					aNo, _ = strconv.ParseUint(no_, 10, 64)
				}

				aPoem := NewPoem(aNo, title_, dynasty_, author_, content_)
				if err := poems.Add(aPoem); err != nil {
					dialog.ShowError(err, win)
					return
				}
				mgr.SwitchToWithCtx("detail", NewDetailContext(aPoem, nil))
			}
		}
	})
	cancelBtn := widget.NewButtonWithIcon("取消", theme.NavigateBackIcon(), func() {
		if ctx, err := context.Get(); err != nil || ctx == nil {
			return
		} else {
			p := ctx.(*EditContext).poem
			if p != nil { // switch from detail, switch back
				mgr.SwitchTo("detail")
			} else {
				mgr.SwitchTo("entry")
			}
		}
	})

	context.AddListener(binding.NewDataListener(func() {
		ctx, err := context.Get()
		if err != nil || ctx == nil {
			return
		}
		p := ctx.(*EditContext).poem
		if p != nil {
			no.Text = fmt.Sprintf("%d", p.No)
			title.Text = p.Title
			author.Text = p.Author
			dynasty.Text = p.Dynasty
			content.Text = p.Content
		} else {
			no.Text = ""
			title.Text = ""
			author.Text = ""
			dynasty.Text = ""
			content.Text = ""
		}

		no.Refresh()
		title.Refresh()
		author.Refresh()
		dynasty.Refresh()
		content.Refresh()
	}))

	return &EditScreen{
		root: container.NewBorder(nil, container.NewGridWithColumns(2, cancelBtn, saveBtn), nil, nil, editable),
		ctx:  context,
	}
}

func (s *EditScreen) Show(ctx interface{}) {
	if ctx == nil {
		ctx = NewEditContext(nil)
	}
	_ = s.ctx.Set(ctx)
	s.root.Show()
}

func (s *EditScreen) Hide() {
	s.root.Hide()
}

func (s *EditScreen) RootObj() fyne.CanvasObject {
	return s.root
}

package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type EntryScreen struct {
	root   fyne.CanvasObject
	update func()
}

func NewEntryScreen(poems *Poems, mgr *ScreenManager, win fyne.Window) *EntryScreen {
	search := binding.NewUntyped()
	_ = search.Set(EmptySearch())

	rule := binding.NewString()
	favorOnly := binding.NewBool()
	updateSearch := func() {
		rule_, _ := rule.Get()
		favorOnly_, _ := favorOnly.Get()
		_ = search.Set(NewSearch(rule_, favorOnly_))
	}
	rule.AddListener(binding.NewDataListener(updateSearch))
	favorOnly.AddListener(binding.NewDataListener(updateSearch))

	favorCheck := widget.NewCheckWithData("仅收藏", favorOnly)
	ruleEntry := widget.NewEntryWithData(rule)
	ruleEntry.SetPlaceHolder("请输入要搜索的词")
	clearRuleBtn := widget.NewButtonWithIcon("清空", theme.ContentClearIcon(), func() {
		ruleEntry.SetText("")
	})
	searchBar := container.NewBorder(nil, nil, favorCheck, clearRuleBtn, ruleEntry)

	poemData := binding.NewUntypedList()
	poemBrowserList := widget.NewListWithData(poemData,
		func() fyne.CanvasObject {
			abstract := widget.NewLabel("")
			showDetailBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			})
			toggleFavorBtn := widget.NewButtonWithIcon("", starOutlineSvg, func() {})
			return container.NewBorder(nil, nil, showDetailBtn, toggleFavorBtn, abstract)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			item.AddListener(binding.NewDataListener(func() {
				i, _ := item.(binding.Untyped).Get()
				p := i.(*Poem)

				objs := o.(*fyne.Container).Objects
				abstract, showDetailBtn, toggleFavorBtn := objs[0].(*widget.Label), objs[1].(*widget.Button), objs[2].(*widget.Button)

				abstract.SetText(p.Abstract())

				showDetailBtn.OnTapped = func() {
					s, _ := search.Get()
					search_ := s.(*Search)
					mgr.SwitchToWithCtx("detail", NewDetailContext(p, search_))
				}

				updateToggleFavorBtn := func(favor bool) {
					if favor {
						toggleFavorBtn.SetIcon(starFillSvg)
					} else {
						toggleFavorBtn.SetIcon(starOutlineSvg)
					}
				}

				updateToggleFavorBtn(p.Favor)

				toggleFavorBtn.OnTapped = func() {
					p.Favor = !p.Favor
					updateToggleFavorBtn(p.Favor)
				}
			}))
		})

	poemSearchList := widget.NewListWithData(poemData,
		func() fyne.CanvasObject {
			preview := widget.NewRichTextWithText("\n")
			showDetailBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			})
			toggleFavorBtn := widget.NewButtonWithIcon("", starOutlineSvg, func() {})
			return container.NewBorder(nil, nil, showDetailBtn, toggleFavorBtn, preview)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			item.AddListener(binding.NewDataListener(func() {
				i, _ := item.(binding.Untyped).Get()
				p := i.(*Poem)

				objs := o.(*fyne.Container).Objects
				preview, showDetailBtn, toggleFavorBtn := objs[0].(*widget.RichText), objs[1].(*widget.Button), objs[2].(*widget.Button)

				s, _ := search.Get()
				search_ := s.(*Search)
				preview.ParseMarkdown(p.PreviewMarkdown(search_))

				showDetailBtn.OnTapped = func() {
					mgr.SwitchToWithCtx("detail", NewDetailContext(p, search_))
				}

				updateToggleFavorBtn := func(favor bool) {
					if favor {
						toggleFavorBtn.SetIcon(starFillSvg)
					} else {
						toggleFavorBtn.SetIcon(starOutlineSvg)
					}
				}

				updateToggleFavorBtn(p.Favor)

				toggleFavorBtn.OnTapped = func() {
					p.Favor = !p.Favor
					updateToggleFavorBtn(p.Favor)
				}
			}))
		})

	showSearchList := func(b bool) {
		if b {
			poemSearchList.Refresh()
			poemSearchList.Show()
			poemBrowserList.Hide()
		} else {
			poemSearchList.Hide()
			poemBrowserList.Refresh()
			poemBrowserList.Show()
		}
	}
	showSearchList(false)

	updateList := func() {
		s, _ := search.Get()
		search_ := s.(*Search)
		filteredPoems := poems.Filter(search_)

		filtered := make([]interface{}, len(filteredPoems))
		for i := range filtered {
			filtered[i] = filteredPoems[i]
		}

		_ = poemData.Set(filtered)
		showSearchList(search_.HasKeyword())
	}

	gotoBtn := widget.NewButtonWithIcon("跳转", theme.SearchIcon(), func() {
		text := widget.NewEntry()
		text.Validator = func(s string) error {
			if _, err := strconv.ParseUint(s, 10, 64); err != nil {
				return errors.New("请输入正整数")
			}
			return nil
		}
		dialog.ShowForm("提示", "跳转", "关闭", []*widget.FormItem{widget.NewFormItem("序号", text)}, func(b bool) {
			if !b {
				return
			}

			id, err := strconv.ParseUint(text.Text, 10, 64)
			if err != nil {
				return
			}

			list, _ := poemData.Get()
			for row, p := range list {
				poem := p.(*Poem)
				if poem.ID == id {
					poemBrowserList.Select(row)
					poemSearchList.Select(row)
					break
				}
			}
		}, win)
	})
	exportBtn := widget.NewButtonWithIcon("导出", theme.DocumentSaveIcon(), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}

			if err := poems.Export(writer); err != nil {
				dialog.ShowError(err, win)
			} else {
				dialog.ShowInformation(" 提示", "导出成功", win)
			}
		}, win)
	})
	importBtn := widget.NewButtonWithIcon("导入", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}

			defer func() {
				err = reader.Close()
			}()
			if err := poems.Import(reader); err != nil {
				dialog.ShowError(err, win)
			}
			updateList()
		}, win)
	})
	addBtn := widget.NewButtonWithIcon("添加", theme.ContentAddIcon(), func() {
		// TODO add poem
	})

	search.AddListener(binding.NewDataListener(updateList))

	root := container.NewBorder(searchBar, container.NewGridWithColumns(4, gotoBtn, exportBtn, importBtn, addBtn), nil, nil, container.NewMax(poemBrowserList, poemSearchList))

	return &EntryScreen{root: root, update: updateList}
}

func (s *EntryScreen) Show(update_ interface{}) {
	if update_ != nil && update_.(bool) {
		s.update()
	}
	s.root.Show()
}

func (s *EntryScreen) Hide() {
	s.root.Hide()
}

func (s *EntryScreen) RootObj() fyne.CanvasObject {
	return s.root
}

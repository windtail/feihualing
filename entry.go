package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EntryScreen struct {
	root fyne.CanvasObject
}

func NewEntryScreen(poems *Poems, mgr *ScreenManager) *EntryScreen {
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
	poemList := widget.NewListWithData(poemData,
		func() fyne.CanvasObject {
			abstract := widget.NewLabel("")
			preview := widget.NewRichTextWithText("")
			showDetailBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			})
			return container.NewBorder(nil, nil, showDetailBtn, nil, container.NewGridWithRows(2, abstract, preview))
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			i, _ := item.(binding.Untyped).Get()
			p := i.(*Poem)

			objs := o.(*fyne.Container).Objects
			grid, showDetailBtn := objs[0].(*fyne.Container), objs[1].(*widget.Button)
			objs = grid.Objects
			abstract, preview := objs[0].(*widget.Label), objs[1].(*widget.RichText)

			abstract.SetText(p.Abstract())
			s, _ := search.Get()
			search_ := s.(*Search)
			preview.ParseMarkdown(p.PreviewMarkdown(search_))

			showDetailBtn.OnTapped = func() {
				mgr.SwitchToWithCtx("detail", NewDetailContext(p, search_))
			}
		})

	addPoemBtn := widget.NewButtonWithIcon("添加", theme.ContentAddIcon(), func() {
		// TODO add poem
	})

	search.AddListener(binding.NewDataListener(func() {
		s, _ := search.Get()
		search_ := s.(*Search)
		filteredPoems := poems.Filter(search_)

		filtered := make([]interface{}, len(filteredPoems))
		for i := range filtered {
			filtered[i] = filteredPoems[i]
		}

		_ = poemData.Set(filtered)
		poemList.Refresh()
	}))

	root := container.NewBorder(searchBar, addPoemBtn, nil, nil, poemList)

	return &EntryScreen{root: root}
}

func (s *EntryScreen) Show(interface{}) {
	s.root.Show()
}

func (s *EntryScreen) Hide() {
	s.root.Hide()
}

func (s *EntryScreen) RootObj() fyne.CanvasObject {
	return s.root
}

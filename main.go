package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
)

//go:embed poems.json
var poemsSrcData []byte

type DetailPoemParam struct {
	poem *Poem
	key  string
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("飞花令")
	myApp.Settings().SetTheme(&myTheme{})
	mgr := NewScreenManager()

	if _, ok := myApp.(desktop.App); ok {
		myWindow.Resize(fyne.NewSize(800, 600))
	}

	// 一个诗句的详情页
	detailPoemParam := binding.NewUntyped()
	detailText := widget.NewRichTextWithText("")
	returnBtn := widget.NewButtonWithIcon("返回", theme.NavigateBackIcon(), func() {})
	delPoemBtn := widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {
		if param, err := detailPoemParam.Get(); err != nil || param == nil {
			return
		} else {
			p := param.(*DetailPoemParam).poem
			dialog.ShowConfirm("警告", fmt.Sprintf("确认要删除 %s 吗？", p.Title), func(b bool) {
				if !b {
					return
				}

			}, myWindow)
		}
	})
	detailScreen := container.NewBorder(nil, container.NewGridWithColumns(2, returnBtn, delPoemBtn), nil, nil, detailText)
	returnBtn.OnTapped = func() {
		mgr.SwitchTo("main")
	}

	detailPoemParam.AddListener(binding.NewDataListener(func() {
		param, err := detailPoemParam.Get()
		if err != nil || param == nil {
			return
		}
		p := param.(*DetailPoemParam)

		detailText.ParseMarkdown(p.poem.Markdown(p.key))
	}))
	mgr.Add("detail", detailScreen)

	// 主页面
	searchText := binding.NewString()
	searchEntry := widget.NewEntryWithData(searchText)
	searchEntry.SetPlaceHolder("请输入要搜索的词")
	searchClearBtn := widget.NewButtonWithIcon("清空", theme.ContentClearIcon(), func() {
		searchEntry.SetText("")
	})
	searchBar := container.NewBorder(nil, nil, nil, searchClearBtn, searchEntry)

	poems := make([]*Poem, 0)
	_ = json.Unmarshal(poemsSrcData, &poems)

	poemData := binding.NewUntypedList()
	var selectedPoemId int64 = -1
	poemList := widget.NewListWithData(poemData,
		func() fyne.CanvasObject {
			abstract := widget.NewLabel("")
			showDetailBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			})
			showDetailBtn.Hide()
			return container.NewBorder(nil, nil, nil, showDetailBtn, abstract)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			i, _ := item.(binding.Untyped).Get()
			p := i.(*Poem)

			objs := o.(*fyne.Container).Objects
			abstract, showDetailBtn := objs[0].(*widget.Label), objs[1].(*widget.Button)
			abstract.SetText(p.Abstract())

			if selectedPoemId == p.Id {
				showDetailBtn.OnTapped = func() {
					key, _ := searchText.Get()
					_ = detailPoemParam.Set(&DetailPoemParam{
						poem: p,
						key:  strings.TrimSpace(key),
					})
					mgr.SwitchTo("detail")
				}
				showDetailBtn.Show()
			} else {
				showDetailBtn.Hide()
			}
		})
	poemList.OnSelected = func(id widget.ListItemID) {
		if v, err := poemData.GetValue(id); err != nil {
			return
		} else {
			poem := v.(*Poem)
			selectedPoemId = poem.Id
			poemList.Refresh()
		}
	}

	addPoemBtn := widget.NewButtonWithIcon("添加", theme.ContentAddIcon(), func() {
	})
	mainScreen := container.NewBorder(searchBar, addPoemBtn, nil, nil, poemList)

	searchText.AddListener(binding.NewDataListener(func() {
		key := strings.TrimSpace(searchEntry.Text)
		var filtered []interface{}

		if len(key) == 0 {
			for _, poem := range poems {
				filtered = append(filtered, poem)
			}
			searchClearBtn.Hide()
		} else {
			for _, poem := range poems {
				if strings.Contains(poem.Content, key) {
					filtered = append(filtered, poem)
				}
			}
			searchClearBtn.Show()
		}

		_ = poemData.Set(filtered)
		selectedPoemId = -1
		poemList.UnselectAll()
	}))
	mgr.Add("main", mainScreen)

	mgr.SwitchTo("main")
	content := container.NewScroll(container.NewMax(mainScreen, detailScreen))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

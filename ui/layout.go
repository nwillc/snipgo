package ui

import (
	"github.com/atotto/clipboard"
	"github.com/nwillc/snipgo/model"
	"github.com/rivo/tview"
)

var (
	rowsWeights = []int{3, 0, 3}
	colWeights  = []int{20, 45, 0, 10}
	mainRow = 1
	footerRow = 2
)

type UI struct {
	app *tview.Application
	*tview.Grid
	editor       *Editor
	categoryList *tview.List
	titleList    *tview.List
	categories   *model.Categories
}

func NewLayout() *UI {
	app := tview.NewApplication()
	editor := NewEditor()
	grid := tview.NewGrid().
		SetRows(rowsWeights...).
		SetColumns(colWeights...).
		SetBorders(true)

	categoryList := tview.NewList().
		ShowSecondaryText(false)

	titleList := tview.NewList().
		ShowSecondaryText(false)

	copyButton := tview.NewButton("Copy").SetSelectedFunc(func() {
		clipboard.WriteAll(editor.String())
	})

	grid.
		AddItem(categoryList, mainRow, 0, 1, 1, 0, 100, true).
		AddItem(titleList, mainRow, 1, 1, 1, 0, 100, true).
		AddItem(editor, mainRow, 2, 1, 2, 0, 100, false).
		AddItem(copyButton, footerRow, 3, 1, 1, 0, 0, true)

	ui := &UI{app, grid, editor, categoryList, titleList, nil}

	categoryList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		ui.loadTitles()
	})

	titleList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		ui.editor.Text(ui.CurrentSnippet().Body)
	})

	return ui
}

func (ui *UI) Categories(categories *model.Categories) {
	ui.categories = categories
	ui.loadCategories()
	ui.loadTitles()
}

func (ui *UI) CurrentCategory() *model.Category {
	selected := ui.categoryList.GetCurrentItem()
	return &(*ui.categories)[selected]
}

func (ui *UI) CurrentSnippet() *model.Snippet {
	selected := ui.titleList.GetCurrentItem()
	return &ui.CurrentCategory().Snippets[selected]
}

func (ui *UI) loadCategories() {
	ui.categoryList.Clear()
	if ui.categories != nil {
		for _, category := range *ui.categories {
			ui.categoryList.AddItem(category.Name, "", 0, nil)
		}
	}
}

func (ui *UI) loadTitles() {
	ui.titleList.Clear()
	for _, snippet := range ui.CurrentCategory().Snippets {
		ui.titleList.AddItem(snippet.Title, "", 0, nil)
	}
}

func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

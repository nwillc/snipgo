package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sort"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	preferences, err := GetPreferences()
	if err != nil {
		panic("Could not get preferences")
	}
	snippets, err := ReadSnippets(preferences.DefaultFile)
	if err != nil {
		panic("Could not read Snippets")
	}

	categories := SnippetsByCategory(snippets)
	sort.Sort(categories)

	app := tview.NewApplication()

	layout := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	titlesTable := tview.NewTable().
		SetBorders(false)

	lastSnippets := &categories[0].Snippets
	loadTitles(titlesTable, lastSnippets)

	categoryTable := tview.NewTable().
		SetBorders(false)
	loadCategories(categoryTable, categories)
	categoryTable.
		Select(0, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				categoryTable.SetSelectable(true, true)
			}
			if key == tcell.KeyRight {
				app.SetFocus(titlesTable)
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			lastSnippets = &categories[row].Snippets
			loadTitles(titlesTable, lastSnippets)
			app.SetFocus(titlesTable)
		})

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	titlesTable.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				titlesTable.SetSelectable(true, true)
			}
			if key == tcell.KeyLeft {
				app.SetFocus(categoryTable)
			}
		}).
		SetSelectedFunc(func(row, column int) {
			textView.SetText((*lastSnippets)[row].Body)
		})

	layout.AddItem(categoryTable, 1, 0, 1, 1, 0, 100, true).
		AddItem(titlesTable, 1, 1, 1, 1, 0, 100, true).
		AddItem(textView, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func loadCategories(t *tview.Table, categories Categories) {
	t.Clear()
	cols, rows := 1, len(categories)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			t.SetCell(r, c,
				tview.NewTableCell(categories[r].Name).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}
}

func loadTitles(t *tview.Table, snippets *Snippets) {
	t.Clear()
	cols, rows := 1, len(*snippets)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			t.SetCell(r, c,
				tview.NewTableCell((*snippets)[r].Title).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}
}

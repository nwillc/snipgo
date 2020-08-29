package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"sort"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not get home directory")
	}
	prefs := fmt.Sprintf("%s/%s", home, preferencesFile)
	preferences, err := ReadPreferences(prefs)
	if err != nil {
		panic("Could not get preferences")
	}
	fmt.Printf("Snippets at: %s\n", preferences.DefaultFile)
	snippets, err := ReadSnippets(preferences.DefaultFile)
	if err != nil {
		panic("Could not read Snippets")
	}
	fmt.Printf("Read %d Snippets\n", len(snippets))

	categories := SnippetsByCategory(snippets)
	sort.Sort(categories)

	app := tview.NewApplication()

	titles := newPrimitive("Titles")
	snippet := newPrimitive("Snippet")

	categoryList := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false)

	table := tview.NewTable().
		SetBorders(false)
	cols, rows := 1, len(categories)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			table.SetCell(r, c,
				tview.NewTableCell(categories[r].Name).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}
	table.
		Select(0, 0).
		SetFixed(1, 1).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			cell := table.GetCell(row, column)
			if cell.Color != tcell.ColorRed {
				cell.SetTextColor(tcell.ColorRed)
			} else {
				cell.SetTextColor(tcell.ColorWhite)
			}
		})

	categoryList.AddItem(table, 1, 0, 1, 1, 0, 100, true).
		AddItem(titles, 1, 1, 1, 1, 0, 100, false).
		AddItem(snippet, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(categoryList, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

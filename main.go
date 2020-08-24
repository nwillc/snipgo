package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"sort"
)

func main() {
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
	sort.Sort(Snippets(snippets))

	distinctCategories := make(map[string]bool)
	categoryList := make([]string, 0)

	for _, snippet := range snippets {
		_, ok := distinctCategories[snippet.Category]
		if !ok {
			distinctCategories[snippet.Category] = true
			categoryList = append(categoryList, snippet.Category)
		}
	}

	sort.Strings(categoryList)

	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false)
	cols, rows := 1, len(categoryList)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			table.SetCell(r, c,
				tview.NewTableCell(categoryList[r]).
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
			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			table.SetSelectable(false, false)
		})
	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

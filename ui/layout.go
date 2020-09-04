/*
 * Copyright (c) 2020, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package ui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/nwillc/snipgo/model"
	"github.com/rivo/tview"
	"sort"
)

var (
	rowsWeights = []int{3, 0, 0, 3}
	colWeights  = []int{25, 0, 0, 10}
	headerRow   = 0
	browserRow  = 1
	editorRow   = 2
	footerRow   = 3
)

type UI struct {
	app *tview.Application
	*tview.Grid
	editor          *Editor
	categoryList    *tview.List
	titleList       *tview.List
	categories      *model.Categories
	currentCategory int
	currentSnippet  int
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

	testButton := tview.NewButton("TEST")

	grid.
		AddItem(categoryList, browserRow, 0, 1, 1, 0, 100, true).
		AddItem(titleList, browserRow, 1, 1, 3, 0, 100, true).
		AddItem(editor, editorRow, 0, 1, 4, 0, 100, false).
		AddItem(copyButton, footerRow, 3, 1, 1, 0, 0, true).
		AddItem(testButton, headerRow, 3, 1, 1, 0, 0, true)

	ui := UI{
		app,
		grid,
		editor,
		categoryList,
		titleList,
		nil,
		-1,
		-1,
	}

	categoryList.SetChangedFunc(func(i int, s string, s2 string, r rune) {
		ui.SetCurrentCategory(i)
	})

	titleList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		ui.SetCurrentSnippet(i)
	})

	testButton.SetSelectedFunc(func() {
		category := model.Category{"TEST", nil}
		added := append(*ui.categories, category)
		sort.Sort(added)
		ui.Categories(&added)
	})

	return &ui
}

func (ui *UI) Categories(categories *model.Categories) {
	ui.categories = categories
	ui.loadCategories()
}

func (ui *UI) CurrentCategory() (*model.Category, error) {
	if ui.currentCategory < 0 || ui.currentCategory >= len(*ui.categories) {
		return nil, fmt.Errorf("no category selected")
	}
	return &(*ui.categories)[ui.currentCategory], nil
}

func (ui *UI) SetCurrentCategory(i int) {
	if i <= len(*ui.categories) {
		ui.currentCategory = i
		ui.loadTitles()
	}
}

func (ui *UI) CurrentSnippet() (*model.Snippet, error) {
	category, err := ui.CurrentCategory()
	if err != nil {
		return nil, err
	}
	if ui.currentSnippet < 0 || ui.currentSnippet >= len(category.Snippets) {
		return nil, fmt.Errorf("no snippet selected")
	}

	return &category.Snippets[ui.currentSnippet], nil
}

func (ui *UI) SetCurrentSnippet(i int) {
	category, err := ui.CurrentCategory()
	if err != nil || len(category.Snippets) == 0 {
		ui.currentSnippet = -1
		ui.editor.Text("")
		return
	}

	ui.currentSnippet = i
	ui.loadSnippet()
}

func (ui *UI) loadCategories() {
	ui.categoryList.Clear()
	if ui.categories != nil {
		for _, category := range *ui.categories {
			ui.categoryList.AddItem(category.Name, "", 0, nil)
		}
	}
	ui.SetCurrentCategory(0)
}

func (ui *UI) loadTitles() {
	ui.titleList.Clear()
	category, err := ui.CurrentCategory()
	if err != nil {
		return
	}
	for _, snippet := range category.Snippets {
		ui.titleList.AddItem(snippet.Title, "", 0, nil)
	}
	ui.SetCurrentSnippet(0)
}

func (ui *UI) loadSnippet() {
	snippet, err := ui.CurrentSnippet()
	if err != nil {
		return
	}
	ui.editor.Text(snippet.Body)
}

func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

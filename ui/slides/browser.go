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

package slides

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/editor"
	"github.com/rivo/tview"
)

var (
	rowsWeights = []int{3, 0, 0, 3}
	colWeights  = []int{25, 0, 0, 10}
	browserRow  = 1
	editorRow   = 2
	footerRow   = 3
)

type BrowserPage struct {
	*tview.Grid
	editor          *editor.Editor
	categoryList    *tview.List
	titleList       *tview.List
	categories      *model.Categories
	currentCategory int
	currentSnippet  int
}

// Implements SetCategories
var _ model.SetCategories = (*BrowserPage)(nil)

func NewBrowserPage() *BrowserPage {
	editor := editor.NewEditor()
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

	page := BrowserPage{
		grid,
		editor,
		categoryList,
		titleList,
		nil,
		-1,
		-1,
	}

	page.
		AddItem(categoryList, browserRow, 0, 1, 1, 0, 100, true).
		AddItem(titleList, browserRow, 1, 1, 3, 0, 100, true).
		AddItem(editor, editorRow, 0, 1, 4, 0, 100, false).
		AddItem(copyButton, footerRow, 3, 1, 1, 0, 0, true)

	categoryList.SetChangedFunc(func(i int, s string, s2 string, r rune) {
		page.setCurrentCategory(i)
	})

	titleList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		page.setCurrentSnippet(i)
	})

	return &page
}

func (ui *BrowserPage) SetCategories(categories *model.Categories) {
	ui.categories = categories
	ui.loadCategories()
}

func (ui *BrowserPage) loadCategories() {
	ui.categoryList.Clear()
	if ui.categories != nil {
		for _, category := range *ui.categories {
			ui.categoryList.AddItem(category.Name, "", 0, nil)
		}
	}
	ui.setCurrentCategory(0)
}

func (browserPage *BrowserPage) setCurrentCategory(i int) {
	if i <= len(*browserPage.categories) {
		browserPage.currentCategory = i
		browserPage.loadTitles()
	}
}

func (ui *BrowserPage) loadTitles() {
	ui.titleList.Clear()
	category, err := ui.getCurrentCategory()
	if err != nil {
		return
	}
	for _, snippet := range category.Snippets {
		ui.titleList.AddItem(snippet.Title, "", 0, nil)
	}
	ui.setCurrentSnippet(0)
}

func (ui *BrowserPage) getCurrentCategory() (*model.Category, error) {
	if ui.currentCategory < 0 || ui.currentCategory >= len(*ui.categories) {
		return nil, fmt.Errorf("no category selected")
	}
	return &(*ui.categories)[ui.currentCategory], nil
}

func (ui *BrowserPage) setCurrentSnippet(i int) {
	category, err := ui.getCurrentCategory()
	if err != nil || len(category.Snippets) == 0 {
		ui.currentSnippet = -1
		ui.editor.Text("")
		return
	}

	ui.currentSnippet = i
	ui.loadSnippet()
}

func (ui *BrowserPage) loadSnippet() {
	snippet, err := ui.getCurrentSnippet()
	if err != nil {
		return
	}
	ui.editor.Text(snippet.Body)
}

func (ui *BrowserPage) getCurrentSnippet() (*model.Snippet, error) {
	category, err := ui.getCurrentCategory()
	if err != nil {
		return nil, err
	}
	if ui.currentSnippet < 0 || ui.currentSnippet >= len(category.Snippets) {
		return nil, fmt.Errorf("no snippet selected")
	}

	return &category.Snippets[ui.currentSnippet], nil
}

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

package pages

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/widgets"
	"github.com/rivo/tview"
)

var (
	rowsWeights = []int{0, 0, 1}
	colWeights  = []int{25, 0}
	browserRow  = 0
	editorRow   = 1
	footerRow   = 2
)

type BrowserPage struct {
	tview.Primitive
	editor          *widgets.Editor
	categoryList    *tview.List
	titleList       *tview.List
	categories      *model.Categories
	currentCategory int
	currentSnippet  int
}

// Implements Slide
var _ Slide = (*BrowserPage)(nil)

func NewBrowserPage() *BrowserPage {
	editor := widgets.NewEditor()
	grid := tview.NewGrid().
		SetRows(rowsWeights...).
		SetColumns(colWeights...).
		SetBorders(true)

	categoryList := tview.NewList().
		ShowSecondaryText(false)

	titleList := tview.NewList().
		ShowSecondaryText(false)

	menu := widgets.NewMenuBar().
		AddItem("Copy", func(i int) {
			clipboard.WriteAll(editor.String())
		}).
		AddItem("Save", func(i int) {

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

	grid.
		AddItem(categoryList, browserRow, 0, 1, 1, 0, 100, true).
		AddItem(titleList, browserRow, 1, 1, 1, 0, 100, true).
		AddItem(editor, editorRow, 0, 1, 2, 0, 100, false).
		AddItem(menu, footerRow, 0, 1, 2, 0, 0, true)

	categoryList.SetChangedFunc(func(i int, s string, s2 string, r rune) {
		page.setCurrentCategory(i)
	})

	titleList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		page.setCurrentSnippet(i)
	})

	return &page
}

func (browserPage *BrowserPage) GetName() string {
	return "Browser"
}

func (browserPage *BrowserPage) SetCategories(categories *model.Categories) {
	browserPage.categories = categories
	browserPage.loadCategories()
}

func (browserPage *BrowserPage) loadCategories() {
	browserPage.categoryList.Clear()
	if browserPage.categories != nil {
		for _, category := range *browserPage.categories {
			browserPage.categoryList.AddItem(category.Name, "", 0, nil)
		}
	}
	browserPage.setCurrentCategory(0)
}

func (browserPage *BrowserPage) setCurrentCategory(i int) {
	if i <= len(*browserPage.categories) {
		browserPage.currentCategory = i
		browserPage.loadTitles()
	}
}

func (browserPage *BrowserPage) loadTitles() {
	browserPage.titleList.Clear()
	category, err := browserPage.getCurrentCategory()
	if err != nil {
		return
	}
	for _, snippet := range category.Snippets {
		browserPage.titleList.AddItem(snippet.Title, "", 0, nil)
	}
	browserPage.setCurrentSnippet(0)
}

func (browserPage *BrowserPage) getCurrentCategory() (*model.Category, error) {
	if browserPage.currentCategory < 0 || browserPage.currentCategory >= len(*browserPage.categories) {
		return nil, fmt.Errorf("no category selected")
	}
	return &(*browserPage.categories)[browserPage.currentCategory], nil
}

func (browserPage *BrowserPage) setCurrentSnippet(i int) {
	category, err := browserPage.getCurrentCategory()
	if err != nil || len(category.Snippets) == 0 {
		browserPage.currentSnippet = -1
		browserPage.editor.Text("")
		return
	}

	browserPage.currentSnippet = i
	browserPage.loadSnippet()
}

func (browserPage *BrowserPage) loadSnippet() {
	snippet, err := browserPage.getCurrentSnippet()
	if err != nil {
		return
	}
	browserPage.editor.Text(snippet.Body)
}

func (browserPage *BrowserPage) getCurrentSnippet() (*model.Snippet, error) {
	category, err := browserPage.getCurrentCategory()
	if err != nil {
		return nil, err
	}
	if browserPage.currentSnippet < 0 || browserPage.currentSnippet >= len(category.Snippets) {
		return nil, fmt.Errorf("no snippet selected")
	}

	return &category.Snippets[browserPage.currentSnippet], nil
}

func (browserPage *BrowserPage) SetCategoryReceiver(receiver CategoryReceiver) {
	// NoOp
}

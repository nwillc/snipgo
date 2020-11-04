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
	"github.com/nwillc/snipgo/services"
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

// BrowserPage is the Slide implementation for browsing model.Categories.
type BrowserPage struct {
	tview.Primitive
	ctx             *services.Context
	editor          *widgets.Editor
	categoryList    *tview.List
	titleList       *tview.List
	categories      *model.Categories
	currentCategory int
	currentSnippet  int
}

// Implements Slide
var _ Slide = (*BrowserPage)(nil)

// NewBrowserPage is a factory for BrowserPage.
func NewBrowserPage(ctx *services.Context) *BrowserPage {
	editor := widgets.NewEditor()
	grid := tview.NewGrid().
		SetRows(rowsWeights...).
		SetColumns(colWeights...).
		SetBorders(true)

	categoryList := tview.NewList().
		ShowSecondaryText(false)

	titleList := tview.NewList().
		ShowSecondaryText(false)

	menu := widgets.NewButtonBar()

	var page = BrowserPage{
		ctx:             ctx,
		Primitive:       grid,
		editor:          editor,
		categoryList:    categoryList,
		titleList:       titleList,
		currentCategory: -1,
		currentSnippet:  -1,
	}
	menu.
		AddButton("Copy", func() {
			clipboard.WriteAll(editor.String())
		}).
		AddButton("Remove", func() {
			page.removeSnippet()
		}).
		AddButton("Save", func() {
			page.write()
		})

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

// GetName returns the name of this Slide.
func (bp *BrowserPage) GetName() string {
	return "Browser"
}

// Categories sets the model.Categories used on the Slide.
func (bp *BrowserPage) Categories(categories *model.Categories) {
	bp.categories = categories
	bp.loadCategories()
}

func (bp *BrowserPage) write() {
	snippet, err := bp.getCurrentSnippet()
	if err == nil {
		snippet.Body = bp.editor.String()
	}
	bp.categories.ToSnippets().WriteSnippets(bp.ctx.JSON, bp.ctx.OS, bp.ctx.IOUTIL, "")
}

func (bp *BrowserPage) loadCategories() {
	bp.categoryList.Clear()
	if bp.categories != nil {
		for _, category := range *bp.categories {
			bp.categoryList.AddItem(category.Name, "", 0, nil)
		}
	}
	bp.setCurrentCategory(0)
}

func (bp *BrowserPage) setCurrentCategory(i int) {
	if i <= len(*bp.categories) {
		bp.currentCategory = i
		bp.loadTitles()
	}
}

func (bp *BrowserPage) loadTitles() {
	bp.titleList.Clear()
	category, err := bp.getCurrentCategory()
	if err != nil {
		return
	}
	for _, snippet := range category.Snippets {
		bp.titleList.AddItem(snippet.Title, "", 0, nil)
	}
	bp.setCurrentSnippet(0)
}

func (bp *BrowserPage) getCurrentCategory() (*model.Category, error) {
	if bp.currentCategory < 0 || bp.currentCategory >= len(*bp.categories) {
		return nil, fmt.Errorf("no category selected")
	}
	return &(*bp.categories)[bp.currentCategory], nil
}

func (bp *BrowserPage) setCurrentSnippet(i int) {
	category, err := bp.getCurrentCategory()
	if err != nil || len(category.Snippets) == 0 {
		bp.currentSnippet = -1
		bp.editor.Text("")
		return
	}

	bp.currentSnippet = i
	bp.loadSnippet()
}

func (bp *BrowserPage) loadSnippet() {
	snippet, err := bp.getCurrentSnippet()
	if err != nil {
		return
	}
	bp.editor.Text(snippet.Body)
}

func (bp *BrowserPage) getCurrentSnippet() (*model.Snippet, error) {
	category, err := bp.getCurrentCategory()
	if err != nil {
		return nil, err
	}
	if bp.currentSnippet < 0 || bp.currentSnippet >= len(category.Snippets) {
		return nil, fmt.Errorf("no snippet selected")
	}

	return &category.Snippets[bp.currentSnippet], nil
}

func (bp *BrowserPage) removeSnippet() {
	target, err := bp.getCurrentSnippet()
	if err != nil {
		return
	}
	for ci, category := range *bp.categories {
		if category.Name != target.Category {
			continue
		}
		for si, snippet := range category.Snippets {
			if snippet.Title == target.Title && snippet.Body == target.Body {
				category.Snippets = append(category.Snippets[:si], category.Snippets[si+1:]...)
				(*bp.categories)[ci] = category
				break
			}
		}
		break
	}
	bp.Categories(bp.categories)
}

// SetCategoryReceiver inform the Slide where to notify with changes of the model.Categories.
func (bp *BrowserPage) SetCategoryReceiver(receiver CategoryReceiver) {
	// NoOp
}

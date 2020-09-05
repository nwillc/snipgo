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
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/pages"
	"github.com/rivo/tview"
	"strconv"
)

type UI struct {
	app *tview.Application
	tview.Primitive
	bp *pages.BrowserPage
}

// Implements SetCategories
var _ model.SetCategories = (*UI)(nil)

func NewUI() *UI {
	app := tview.NewApplication()

	pageView := tview.NewPages()

	_, browserPage := pages.NewBrowserPage()
	aboutPage := pages.NewAboutPage()

	pageNames := []string{"Browser", "About"}

	pageView.
		AddPage("Browser", browserPage, true, true).
		AddPage("About", aboutPage, true, false)

	menu := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pageNo, _ := strconv.Atoi(added[0])
			pageView.SwitchToPage(pageNames[pageNo])
		})

	fmt.Fprintf(menu, `["%d"][darkcyan]%s[white][""]  `, 0, "Browser")
	fmt.Fprintf(menu, `["%d"][darkcyan]%s[white][""]  `, 1, "About")
	menu.Highlight("0")

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(menu, 1, 1, false).
		AddItem(pageView, 0, 1, true)

	ui := UI{
		app,
		layout,
		browserPage,
	}

	return &ui
}

func (ui *UI) SetCategories(categories *model.Categories) {
	ui.bp.SetCategories(categories)
}

func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

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
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/slides"
	"github.com/rivo/tview"
)

type UI struct {
	app *tview.Application
	bp  *slides.BrowserPage
}

// Implements SetCategories
var _ model.SetCategories = (*UI)(nil)

func NewUI() *UI {
	app := tview.NewApplication()
	slide := slides.NewBrowserPage()

	ui := UI{
		app,
		slide,
	}

	return &ui
}

func (ui *UI) SetCategories(categories *model.Categories) {
	ui.bp.SetCategories(categories)
}

func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui.bp, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

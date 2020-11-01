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
	"github.com/nwillc/snipgo/services"
	"github.com/nwillc/snipgo/ui/pages"
	"github.com/nwillc/snipgo/ui/widgets"
	"github.com/rivo/tview"
)

// UI is the tview.Primitive associated with the overall UI.
type UI struct {
	app *tview.Application
	tview.Primitive
	slides []pages.Slide
	pv     *tview.Pages
}

// Implements SetCategories
var _ model.SetCategories = (*UI)(nil)

// NewUI is a factory for UI.
func NewUI(ctx *services.Context) *UI {
	var app = tview.NewApplication()
	var slides = []pages.Slide{
		pages.NewBrowserPage(ctx),
		pages.NewSnippetPage(),
		pages.NewPreferencesPage(),
		pages.NewAboutPage(),
	}
	var pageView = tview.NewPages()
	var menu = widgets.NewButtonBar()
	for i, slide := range slides {
		slideIndex := i
		pageView.AddPage(slide.GetName(), slide, true, i == 0)
		menu.AddButton(slide.GetName(), func() {
			pageView.SwitchToPage(slides[slideIndex].GetName())
		})
	}

	var layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(menu, 1, 1, false).
		AddItem(pageView, 0, 1, true)
	var ui = UI{
		app:       app,
		Primitive: layout,
		slides:    slides,
		pv:        pageView,
	}
	for _, slide := range slides {
		slide.SetCategoryReceiver(func(categories *model.Categories) {
			ui.Categories(categories)
		})
	}

	menu.AddButton("Quit", func() {
		app.Stop()
	})
	return &ui
}

// Categories sets the model.Categories used on the UI.
func (ui *UI) Categories(categories *model.Categories) {
	for _, slide := range ui.slides {
		slide.Categories(categories)
	}
	ui.pv.SwitchToPage("Browser")
}

// Run the UI.
func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

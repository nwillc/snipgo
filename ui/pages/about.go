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
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/widgets"
	"github.com/nwillc/snipgo/gen/version"
	"github.com/rivo/tview"
)

// AboutPage is a Slide showing information about the application.
type AboutPage struct {
	tview.Primitive
}

// Implements Slide
var _ Slide = (*AboutPage)(nil)

// NewAboutPage is factory for AboutPage.
func NewAboutPage() *AboutPage {
	textView := tview.NewTextView()
	fmt.Fprintln(textView, "Snippets Manager")
	fmt.Fprintln(textView, "See https://github.com/nwillc/snipgo")
	fmt.Fprintln(textView, "Version ", version.Version)
	page := AboutPage{widgets.Center(40, 3, textView)}
	return &page
}

// GetName returns the name of this Slide.
func (a *AboutPage) GetName() string {
	return "About"
}

// SetCategories sets the model.Categories used on the Slide.
func (a *AboutPage) SetCategories(categories *model.Categories) {
	// NoOp
}

// SetCategoryReceiver inform the Slide where to notify with changes of the model.Categories.
func (a *AboutPage) SetCategoryReceiver(receiver CategoryReceiver) {
	// NoOp
}

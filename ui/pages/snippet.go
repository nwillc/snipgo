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
	"github.com/nwillc/snipgo/model"
	"github.com/nwillc/snipgo/ui/widgets"
	"github.com/rivo/tview"
	"sort"
)

// SnippetPage is a Slide to add new model.Snippet to the model.Categories.
type SnippetPage struct {
	tview.Primitive
	categories       *model.Categories
	categoryReceiver CategoryReceiver
	category         string
	title            string
	body             string
}

// Implements Slide
var _ Slide = (*SnippetPage)(nil)

// NewSnippetPage is a factory for SnippetPage.
func NewSnippetPage() *SnippetPage {
	var form = tview.NewForm()
	form.SetBorder(true).SetTitle("New Snippet").SetTitleAlign(tview.AlignCenter)
	var page = SnippetPage{
		Primitive: widgets.Center(50, 11, form),
	}
	form.
		AddInputField("Category", "", 20, nil, func(text string) {
			page.setCategory(text)
		}).
		AddInputField("Title", "", 40, nil, func(text string) {
			page.setTitle(text)
		}).
		AddInputField("Body", "", 40, nil, func(text string) {
			page.setBody(text)
		})

	form.AddButton("Save", func() {
		page.addSnippet()
	})
	return &page
}

// SetCategories sets the model.Categories used on the Slide.
func (s *SnippetPage) SetCategories(categories *model.Categories) {
	s.categories = categories
}

// GetName returns the name of this Slide.
func (s *SnippetPage) GetName() string {
	return "Snippet"
}

// SetCategoryReceiver inform the Slide where to notify with changes of the model.Categories.
func (s *SnippetPage) SetCategoryReceiver(receiver CategoryReceiver) {
	s.categoryReceiver = receiver
}

func (s *SnippetPage) addSnippet() {
	updated := *s.categories
	snippet := model.Snippets{{s.category, s.title, s.body}}
	found := false
	for i, category := range updated {
		if category.Name == s.category {
			found = true
			updated[i].Snippets = append(updated[i].Snippets, snippet...)
			break
		}
	}
	if !found {
		category := model.Category{Name: s.category, Snippets: snippet}
		updated = append(updated, category)
	}
	sort.Sort(updated)
	s.categoryReceiver(&updated)
}

func (s *SnippetPage) setCategory(category string) {
	s.category = category
}

func (s *SnippetPage) setTitle(title string) {
	s.title = title
}

func (s *SnippetPage) setBody(body string) {
	s.body = body
}

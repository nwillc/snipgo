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

package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/nwillc/snipgo/snippets"
	"github.com/nwillc/snipgo/ui"
	"github.com/rivo/tview"
)

func main() {
	preferences, err := GetPreferences()
	if err != nil {
		panic("Could not get preferences")
	}
	read, err := snippets.ReadSnippets(preferences.DefaultFile)
	if err != nil {
		panic(fmt.Sprintf("failed %v", err))
	}

	categories := snippets.ByCategory(read)
	lastSnippets := &categories[0].Snippets

	app := tview.NewApplication()

	editor := ui.NewEditor()
	textView := tview.NewTextView()

	titleList := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(i int, s string, s2 string, r rune) {
			editor.Text((*lastSnippets)[i].Body)
		})

	loadTitles(titleList, lastSnippets)

	categoryList := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(i int, s string, s2 string, r rune) {
			lastSnippets = &categories[i].Snippets
			loadTitles(titleList, lastSnippets)
		})

	loadCategories(categoryList, categories)

	copyButton := tview.NewButton("Copy").SetSelectedFunc(func() {
		//if err := clipboard.WriteAll(textView.GetText(false)); err != nil {
		//	panic(fmt.Sprintf("failed to copy to clipboard %v", err))
		//}
		clipboard.WriteAll(editor.String())
	})

	layout := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(textView, 0, 0, 1, 3, 0, 0, true).
		AddItem(copyButton, 2, 0, 1, 3, 0, 0, true).
		AddItem(categoryList, 1, 0, 1, 1, 0, 100, true).
		AddItem(titleList, 1, 1, 1, 1, 0, 100, true).
		AddItem(editor, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func loadCategories(t *tview.List, categories snippets.Categories) {
	t.Clear()
	for _, category := range categories {
		t.AddItem(category.Name, "", 0, nil)
	}
}

func loadTitles(t *tview.List, snippets *snippets.Snippets) {
	t.Clear()
	for _, snippet := range *snippets {
		t.AddItem(snippet.Title, "", 0, nil)
	}
}

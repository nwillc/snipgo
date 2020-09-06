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

package widgets

import (
	"fmt"
	"github.com/rivo/tview"
	"strconv"
)

type MenuBar struct {
	*tview.TextView
	itemCount int
	actions   []func(int)
}

func NewMenuBar() *MenuBar {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	mb := MenuBar{tv, 0, nil}
	mb.SetHighlightedFunc(func(added, removed, remaining []string) {
		index, _ := strconv.Atoi(added[0])
		mb.action(index)
	})
	return &mb
}

func (mb *MenuBar) AddItem(name string, action func(int)) {
	fmt.Fprintf(mb, `|["%d"][darkcyan]%s[white][""]|  `, mb.itemCount, name)
	mb.actions = append(mb.actions, action)
	mb.itemCount += 1
}

func (mb *MenuBar) action(i int) {
	mb.actions[i](i)
}

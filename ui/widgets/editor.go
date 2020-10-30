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
	"github.com/pgavlin/femto"
	"strings"
)

// Editor is a femto.View biased away from files a bit.
type Editor struct {
	buffer *femto.Buffer
	*femto.View
}

// Implements fmt.Stringer
var _ fmt.Stringer = (*Editor)(nil)

// NewEditor is factory function for Editor.
func NewEditor() *Editor {
	var buffer = femto.NewBufferFromString("", "")
	buffer.Settings["ruler"] = false
	var view = femto.NewView(buffer)
	return &Editor{buffer, view}
}

// Text sets the text of the Editor as a string.
func (editor *Editor) Text(text string) {
	editor.buffer.Remove(editor.buffer.Start(), editor.buffer.End())
	editor.buffer.Insert(editor.buffer.Start(), text)
}

// Implement fmt.Stringer
func (editor *Editor) String() string {
	lines := editor.buffer.Lines(0, editor.buffer.NumLines)
	return strings.Join(lines, "\n") + "\n"
}

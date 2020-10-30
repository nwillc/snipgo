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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EditorTestSuite struct {
	suite.Suite
}

func TestEditorTestSuite(t *testing.T) {
	suite.Run(t, new(EditorTestSuite))
}

func (suite *EditorTestSuite) TestNew() {
	editor := NewEditor()
	assert.NotNil(suite.T(), editor)
	assert.Equal(suite.T(), "\n", editor.String())
}

func (suite *EditorTestSuite) TestText() {
	editor := NewEditor()
	text := "This is a\nTest"
	editor.Text(text)
	assert.Equal(suite.T(), text+"\n", editor.String())
}

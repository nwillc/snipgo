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

package model

import (
	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type SnippetsTestSuite struct {
	suite.Suite
	snippets     Snippets
	testFilesDir string
	badFilename  string
	goodFilename string
}

func TestSnippetsTestSuite(t *testing.T) {
	suite.Run(t, new(SnippetsTestSuite))
}

func (suite *SnippetsTestSuite) SetupTest() {
	suite.snippets = []Snippet{
		{
			Category: "A",
			Title:    "A",
		},
		{
			Category: "A",
			Title:    "B",
		},
	}
	suite.testFilesDir = "../test/files"
	suite.badFilename = suite.testFilesDir + "/foo"
	suite.goodFilename = suite.testFilesDir + "/snippets.json"
}

func (suite *SnippetsTestSuite) TestStringer() {
	snippet := Snippet{
		Category: "Foo",
		Title:    "Bar",
		Body:     "Baz",
	}
	assert.Equal(suite.T(), "Foo: Bar", snippet.String())
}

func (suite *SnippetsTestSuite) TestNonExist() {
	_, ok := ReadSnippets(suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestExist() {
	_, ok := ReadSnippets(suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestNoHomeDir() {
	defer monkey.UnpatchAll()

	PatchNoHomeDir()

	defer func() { recover() }()
	_, _ = ReadSnippets("")
	suite.T().Errorf("did not panic")
}

func (suite *SnippetsTestSuite) TestHomeDir() {
	defer monkey.UnpatchAll()

	PatchHomeDir(testFilesDir)

	_, err := ReadSnippets("")
	assert.NoError(suite.T(), err)
}

func (suite *SnippetsTestSuite) TestHomeDirNoPreferences() {
	defer monkey.UnpatchAll()

	PatchHomeDir(".")
	PatchOpenFail()

	_, err := ReadSnippets("")
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "open fail")
	}
}

func (suite *SnippetsTestSuite) TestUnableToRead() {
	defer monkey.UnpatchAll()

	PatchReadAllFail()

	_, err := ReadSnippets(suite.goodFilename)
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "read all fail")
	}
}

func (suite *SnippetsTestSuite) TestWriteDefaultNoHome() {
	defer monkey.UnpatchAll()

	PatchNoHomeDir()

	defer func() { recover() }()
	var snippets = Snippets{}
	_ = snippets.WriteSnippets("")
	suite.T().Errorf("did not panic")
}

func (suite *SnippetsTestSuite) TestWriteDefaultNoPreferences() {
	defer monkey.UnpatchAll()

	PatchHomeDir(".")
	PatchOpenFail()

	var snippets = Snippets{}
	err := snippets.WriteSnippets("")
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "open fail")
	}
}

func (suite *SnippetsTestSuite) TestWriteDefault() {
	defer monkey.UnpatchAll()

	PatchHomeDir(testFilesDir)

	tempFile, err := os.CreateTemp("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	testSnippets, ok := ReadSnippets(suite.goodFilename)
	assert.Nil(suite.T(), ok)

	PatchWriteFileRedirect(tempFile.Name())

	err = testSnippets.WriteSnippets("")
	assert.Nil(suite.T(), err)

	read, ok := ReadSnippets(tempFile.Name())
	assert.Nil(suite.T(), ok)
	assert.Equal(suite.T(), len(testSnippets), len(read))
}

func (suite *SnippetsTestSuite) TestReadUnmarshalFail() {
	defer monkey.UnpatchAll()

	PatchJSONUnmarshalFail()

	_, err := ReadSnippets(suite.goodFilename)
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "json unmarshal fail")
	}
}

func (suite *SnippetsTestSuite) TestWriteMarshalFail() {
	defer monkey.UnpatchAll()

	PatchJSONMarshalFail()

	tempFile, err := os.CreateTemp("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	var testSnippets = Snippets{}
	err = testSnippets.WriteSnippets(suite.goodFilename)
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "json marshal fail")
	}
}

func (suite *SnippetsTestSuite) TestWriteFail() {
	defer monkey.UnpatchAll()
	PatchWriteFileFail()

	tempFile, err := os.CreateTemp("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	testSnippets, ok := ReadSnippets(suite.goodFilename)
	assert.Nil(suite.T(), ok)

	err = testSnippets.WriteSnippets(tempFile.Name())
	if assert.Error(suite.T(), err) {
		assert.Errorf(suite.T(), err, "write file fail")
	}
}

func (suite *SnippetsTestSuite) TestWriteFile() {
	tempFile, err := os.CreateTemp("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	original, ok := ReadSnippets(suite.goodFilename)
	assert.Nil(suite.T(), ok)
	err = original.WriteSnippets(tempFile.Name())
	assert.Nil(suite.T(), err)
	read, ok := ReadSnippets(tempFile.Name())
	assert.Nil(suite.T(), ok)
	assert.Equal(suite.T(), len(original), len(read))
}

func (suite *SnippetsTestSuite) TestMalformedFile() {
	tempFile, err := os.CreateTemp("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString("not json")

	_, ok := ReadSnippets(tempFile.Name())
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestLen() {
	assert.Equal(suite.T(), 2, suite.snippets.Len())
}

func (suite *SnippetsTestSuite) TestLess() {
	assert.LessOrEqual(suite.T(), suite.snippets[0].Category, suite.snippets[1].Category)
	assert.LessOrEqual(suite.T(), suite.snippets[0].Title, suite.snippets[1].Title)
	assert.True(suite.T(), suite.snippets.Less(0, 1))
}

func (suite *SnippetsTestSuite) TestSwap() {
	suite.snippets.Swap(0, 1)
	assert.False(suite.T(), suite.snippets.Less(0, 1))
}

func (suite *SnippetsTestSuite) TestSnippetsByCategory() {
	tests := []struct {
		name       string
		snippets   Snippets
		categories Categories
	}{
		{
			name:       "Empty",
			snippets:   Snippets{},
			categories: Categories{},
		},
		{
			name: "TwoSnippetOneCategory",
			snippets: Snippets{
				{
					Category: "A",
					Title:    "A",
				},
				{
					Category: "A",
					Title:    "B",
				},
			},
			categories: Categories{
				{
					Name:     "A",
					Snippets: []Snippet{{Category: "A", Title: "A"}, {Category: "A", Title: "B"}},
				},
			},
		},
		{
			name: "TwoSnippetTwoCategory",
			snippets: Snippets{
				{
					Category: "A",
					Title:    "A",
				},
				{
					Category: "B",
					Title:    "B",
				},
			},
			categories: Categories{
				{
					Name:     "A",
					Snippets: []Snippet{{Category: "A", Title: "A"}},
				},
				{
					Name:     "B",
					Snippets: []Snippet{{Category: "B", Title: "B"}},
				},
			},
		},
	}
	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			result := test.snippets.ByCategory()
			assert.Equal(t, len(test.categories), len(result))
			for i, category := range test.categories {
				assert.Equal(t, category.Name, result[i].Name)
				assert.Equal(t, len(category.Snippets), len(result[i].Snippets))
				if len(category.Snippets) == len(result[i].Snippets) {
					for j, snippet := range category.Snippets {
						assert.Equal(t, snippet.Title, result[i].Snippets[j].Title)
					}
				}
			}
		})
	}
}

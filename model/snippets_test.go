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
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nwillc/snipgo/mocks"
	"github.com/nwillc/snipgo/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type SnippetsTestSuite struct {
	suite.Suite
	ctx          *services.Context
	snippets     Snippets
	testFilesDir string
	badFilename  string
	goodFilename string
}

func TestSnippetsTestSuite(t *testing.T) {
	suite.Run(t, new(SnippetsTestSuite))
}

func (suite *SnippetsTestSuite) SetupTest() {
	suite.ctx = services.NewDefaultContext()
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
	_, ok := ReadSnippets(suite.ctx, suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestExist() {
	_, ok := ReadSnippets(suite.ctx, suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestNoHomeDir() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return("", fmt.Errorf("foo")).
		Times(1)
	defer func() { recover() }()
	_, _ = ReadSnippets(suite.ctx.CopyUpdateOs(mockOs), "")
	suite.T().Errorf("did not panic")
}

func (suite *SnippetsTestSuite)  TestHomeDir() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
			UserHomeDir().
			Return(testFilesDir, nil).
			Times(1)
	mockOs.EXPECT().
			Open(testFilesDir + "/.snippets.json").
			Return(suite.ctx.OS.Open(suite.testFilesDir + "/preferences.json")).
			Times(1)
	mockOs.EXPECT().
				Open("/Users/nwillc/Documents/snippets.json").
				Return(suite.ctx.OS.Open(suite.goodFilename)).
				Times(1)
	_, err := ReadSnippets(suite.ctx.CopyUpdateOs(mockOs), "")
	assert.Nil(suite.T(), err)
}
func (suite *SnippetsTestSuite) TestWriteMarshalFail() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockJson = mocks.NewMockJson(mockCtrl)
	mockJson.EXPECT().
		Unmarshal(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("mock error")).
		Times(1)
	_, ok := ReadSnippets(suite.ctx.CopyUpdateJson(mockJson), suite.goodFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestWriteFile() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	original, ok := ReadSnippets(suite.ctx, suite.goodFilename)
	assert.Nil(suite.T(), ok)
	err = original.WriteSnippets(suite.ctx, tempFile.Name())
	assert.Nil(suite.T(), err)
	read, ok := ReadSnippets(suite.ctx, tempFile.Name())
	assert.Nil(suite.T(), ok)
	assert.Equal(suite.T(), len(original), len(read))
}

func (suite *SnippetsTestSuite) TestMalformedFile() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString("not json")

	_, ok := ReadSnippets(suite.ctx, tempFile.Name())
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

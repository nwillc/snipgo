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
	snippets     Snippets
	testFilesDir string
	badFilename  string
	goodFilename string
	json         services.Json
	os           services.Os
	ioUtil       services.IoUtil
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
	suite.json = services.NewJson()
	suite.os = services.NewOs()
	suite.ioUtil = services.NewIoUtil()
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
	_, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestExist() {
	_, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, suite.goodFilename)
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
	_, _ = ReadSnippets(suite.json, mockOs, suite.ioUtil, "")
	suite.T().Errorf("did not panic")
}

func (suite *SnippetsTestSuite) TestHomeDir() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return(testFilesDir, nil).
		Times(1)
	mockOs.EXPECT().
		Open(testFilesDir + "/.snippets.json").
		Return(suite.os.Open(suite.testFilesDir + "/preferences.json")).
		Times(1)
	mockOs.EXPECT().
		Open(suite.goodFilename).
		Return(suite.os.Open(suite.goodFilename)).
		Times(1)
	_, err := ReadSnippets(suite.json, mockOs, suite.ioUtil, "")
	assert.Nil(suite.T(), err)
}

func (suite *SnippetsTestSuite) TestHomeDirNoPreferences() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return(testFilesDir, nil).
		Times(1)
	mockOs.EXPECT().
		Open(testFilesDir+"/.snippets.json").
		Return(nil, fmt.Errorf("file not found")).
		Times(1)
	_, err := ReadSnippets(suite.json, mockOs, suite.ioUtil, "")
	assert.NotNil(suite.T(), err)
}

func (suite *SnippetsTestSuite) TestUnableToRead() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockIoUtil = mocks.NewMockIoUtil(mockCtrl)
	mockIoUtil.EXPECT().
		ReadAll(gomock.Any()).
		Return(nil, fmt.Errorf("unable to read")).
		Times(1)
	_, err := ReadSnippets(suite.json, suite.os, mockIoUtil, suite.goodFilename)
	assert.NotNil(suite.T(), err)
	assert.Errorf(suite.T(), err, "unable to read")
}

func (suite *SnippetsTestSuite) TestWriteDefaultNoHome() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return("", fmt.Errorf("foo")).
		Times(1)
	defer func() { recover() }()
	var snippets = Snippets{}
	_ = snippets.WriteSnippets(suite.json, mockOs, suite.ioUtil, "")
	suite.T().Errorf("did not panic")
}

func (suite *SnippetsTestSuite) TestWriteDefaultNoPreferences() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return(testFilesDir, nil).
		Times(1)
	mockOs.EXPECT().
		Open(testFilesDir+"/.snippets.json").
		Return(nil, fmt.Errorf("file not found")).
		Times(1)
	var snippets = Snippets{}
	err := snippets.WriteSnippets(suite.json, mockOs, suite.ioUtil, "")
	assert.NotNil(suite.T(), err)
}

func (suite *SnippetsTestSuite) TestWriteDefault() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	testSnippets, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, suite.goodFilename)
	assert.Nil(suite.T(), ok)
	testSnippetsBytes, err := suite.json.Marshal(testSnippets)
	assert.Nil(suite.T(), err)

	testPrefs, err := ReadPreferences(suite.json, suite.os, suite.ioUtil, suite.testFilesDir+"/preferences.json")
	assert.Nil(suite.T(), err)
	testPrefsBytes, err := suite.json.Marshal(testPrefs)
	assert.Nil(suite.T(), err)

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return(testFilesDir, nil).
		Times(1)
	mockOs.EXPECT().
		Open(testFilesDir + "/.snippets.json").
		Return(suite.os.Open(suite.testFilesDir + "/preferences.json")).
		Times(1)

	var mockIoUtil = mocks.NewMockIoUtil(mockCtrl)
	mockIoUtil.EXPECT().
		ReadAll(gomock.Any()).
		Return(testPrefsBytes, nil).
		Times(1)

	mockIoUtil.EXPECT().
		WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(suite.ioUtil.WriteFile(tempFile.Name(), testSnippetsBytes, os.ModePerm)).
		Times(1)

	err = testSnippets.WriteSnippets(suite.json, mockOs, mockIoUtil, "")
	assert.Nil(suite.T(), err)

	read, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, tempFile.Name())
	assert.Nil(suite.T(), ok)
	assert.Equal(suite.T(), len(testSnippets), len(read))
}

func (suite *SnippetsTestSuite) TestReadMarshalFail() {
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
	_, ok := ReadSnippets(mockJson, suite.os, suite.ioUtil, suite.goodFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestWriteMarshalFail() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockJson = mocks.NewMockJson(mockCtrl)
	mockJson.EXPECT().
		Marshal(gomock.Any()).
		Return(nil, fmt.Errorf("mock error")).
		Times(1)
	var testSnippets = Snippets{}
	ok := testSnippets.WriteSnippets(mockJson, suite.os, suite.ioUtil, suite.goodFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestWriteFail() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	testSnippets, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, suite.goodFilename)
	assert.Nil(suite.T(), ok)

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockIoUtil = mocks.NewMockIoUtil(mockCtrl)
	mockIoUtil.EXPECT().
		WriteFile(tempFile.Name(), gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("foo")).
		Times(1)

	ok = testSnippets.WriteSnippets(suite.json, suite.os, mockIoUtil, tempFile.Name())
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestWriteFile() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	original, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, suite.goodFilename)
	assert.Nil(suite.T(), ok)
	err = original.WriteSnippets(suite.json, suite.os, suite.ioUtil, tempFile.Name())
	assert.Nil(suite.T(), err)
	read, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, tempFile.Name())
	assert.Nil(suite.T(), ok)
	assert.Equal(suite.T(), len(original), len(read))
}

func (suite *SnippetsTestSuite) TestMalformedFile() {
	tempFile, err := ioutil.TempFile("", "snippets.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString("not json")

	_, ok := ReadSnippets(suite.json, suite.os, suite.ioUtil, tempFile.Name())
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

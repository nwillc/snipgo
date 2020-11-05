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

const (
	testFilesDir = "../test/files"
	testPrefFile = testFilesDir + "/preferences.json"
)

type PreferencesTestSuite struct {
	suite.Suite
	badFilename  string
	goodFilename string
	json         services.Json
	os           services.Os
	ioUtil       services.IoUtil
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.badFilename = "foo"
	suite.goodFilename = testPrefFile
	suite.json = services.NewJson()
	suite.os = services.NewOs()
	suite.ioUtil = services.NewIoUtil()
}

func (suite *PreferencesTestSuite) TestNonExistPrefs() {
	_, ok := ReadPreferences(suite.json, suite.os, suite.ioUtil, suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestExistPrefs() {
	_, ok := ReadPreferences(suite.json, suite.os, suite.ioUtil, suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestNoHomeDir() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return("", fmt.Errorf("foo")).
		Times(1)

	defer func() { recover() }()
	_, _ = ReadPreferences(suite.json, mockOs, suite.ioUtil, "")
	suite.T().Errorf("did not panic")
}

func (suite *PreferencesTestSuite) TestHomeDir() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockOs = mocks.NewMockOs(mockCtrl)
	mockOs.EXPECT().
		UserHomeDir().
		Return(testFilesDir, nil).
		Times(1)
	mockOs.EXPECT().
		Open(testFilesDir + "/.snippets.json").
		Return(suite.os.Open(testFilesDir + "/snippets.json")).
		Times(1)
	err, _ := ReadPreferences(suite.json, mockOs, suite.ioUtil, "")
	assert.Nil(suite.T(), err)
}

func (suite *PreferencesTestSuite) TestMalformedFile() {
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.WriteString("not json")

	_, ok := ReadPreferences(suite.json, suite.os, suite.ioUtil, tempFile.Name())
	if assert.NotNil(suite.T(), ok) {
		assert.Errorf(suite.T(), ok, "json marshal failed")
	}
}

func (suite *PreferencesTestSuite) TestWrite() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	err = p.Write(suite.json, suite.ioUtil, tempFile.Name())
	assert.Nil(suite.T(), err)
	read, err := ReadPreferences(suite.json, suite.os, suite.ioUtil, tempFile.Name())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), p.DefaultFile, read.DefaultFile)
}

func (suite *PreferencesTestSuite) TestWriteMarshalFail() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockJson = mocks.NewMockJson(mockCtrl)
	var errMsg = "json marshal failed"
	mockJson.EXPECT().
		Marshal(gomock.Any()).
		Return([]byte{}, fmt.Errorf(errMsg)).
		Times(1)

	err = p.Write(mockJson, suite.ioUtil, tempFile.Name())
	assert.NotNil(suite.T(), err)
	assert.Errorf(suite.T(), err, errMsg)
}

func (suite *PreferencesTestSuite) TestWriteWriteFail() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	var mockIoUtil = mocks.NewMockIoUtil(mockCtrl)
	var errMsg = "write file failed"
	mockIoUtil.EXPECT().
		WriteFile(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fmt.Errorf(errMsg)).
		Times(1)

	err = p.Write(suite.json, mockIoUtil, tempFile.Name())
	assert.NotNil(suite.T(), err)
	assert.Errorf(suite.T(), err, errMsg)
}

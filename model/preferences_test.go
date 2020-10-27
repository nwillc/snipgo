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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

const testPrefFile = "../test/files/preferences.json"

type failingHomeDirMock struct{}

func (h failingHomeDirMock) userHomeDir() (string, error) {
	return "", fmt.Errorf("boom")
}

type PreferencesTestSuite struct {
	suite.Suite
	badFilename  string
	goodFilename string
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.badFilename = "foo"
	suite.goodFilename = testPrefFile
}

func (suite *PreferencesTestSuite) TestBadHomeDir() {
	defaultUserHomeGet := userHomeGet
	userHomeGet = failingHomeDirMock{}
	defer func() {
		userHomeGet = defaultUserHomeGet
		recover()
	}()
	_, _ = ReadPreferences("")
	suite.T().Errorf("expected panic")
}

func (suite *PreferencesTestSuite) TestNonExistPrefs() {
	_, ok := ReadPreferences(suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestExistPrefs() {
	_, ok := ReadPreferences(suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestMalformedFile() {
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	tempFile.WriteString("not json")

	_, ok := ReadPreferences(tempFile.Name())
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestWrite() {
	p := Preferences{DefaultFile: "foo"}
	tempFile, err := ioutil.TempFile("", "prefs.*.json")
	assert.Nil(suite.T(), err)
	defer os.Remove(tempFile.Name())
	err = p.Write(tempFile.Name())
	assert.Nil(suite.T(), err)
	read, err := ReadPreferences(tempFile.Name())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), p.DefaultFile, read.DefaultFile)
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

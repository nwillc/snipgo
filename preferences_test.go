package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const testPrefFile = "test/files/preferences.json"

type PreferencesTestSuite struct {
	suite.Suite
	badFilename string
	goodFilename string
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.badFilename = "foo"
	suite.goodFilename = testPrefFile
}

func (suite *PreferencesTestSuite) TestNonExistPrefs() {
	_, ok := getPreferences(suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestExistPrefs() {
	_, ok := getPreferences(suite.goodFilename)
	assert.Nil(suite.T(), ok)
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

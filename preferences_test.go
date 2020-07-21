package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PreferencesTestSuite struct {
	suite.Suite
	badFilename string
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.badFilename = "foo"
}

func (suite *PreferencesTestSuite) TestDoesNotExist() {
	_, ok := getPreferences(suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *PreferencesTestSuite) TestExist() {
	_, ok := getPreferences(preferencesFile)
	assert.Nil(suite.T(), ok)
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

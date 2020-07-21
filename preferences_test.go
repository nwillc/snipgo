package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PreferencesTestSuite struct {
	suite.Suite
	fileName string
}

func (suite *PreferencesTestSuite) SetupTest() {
	suite.fileName = "foo"
}

func (suite *PreferencesTestSuite) TestDoesNotExist() {
	assert.False(suite.T(), true)
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}

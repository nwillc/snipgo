package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
}

func (suite *MainTestSuite) TestTagExists() {
	repo := GetRepository("../..")
	found := TagExists(repo, "v0.1.0")
	assert.True(suite.T(), found)
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const testSnippetsFile = "test/files/snippets.json"

type SnippetsTestSuite struct {
	suite.Suite
	snippets ByCategoryTitle
	badFilename string
	goodFilename string
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
	suite.badFilename = "foo"
	suite.goodFilename = testSnippetsFile
}

func (suite *SnippetsTestSuite) TestNonExist() {
	_, ok := ReadSnippets(suite.badFilename)
	assert.NotNil(suite.T(), ok)
}

func (suite *SnippetsTestSuite) TestExist() {
	_, ok := ReadSnippets(suite.goodFilename)
	assert.Nil(suite.T(), ok)
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

func TestSnippetsTestSuite(t *testing.T) {
	suite.Run(t, new(SnippetsTestSuite))
}

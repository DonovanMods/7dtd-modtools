package gamexml_test

import (
	"testing"

	"github.com/donovanmods/7dmt/lib/gamexml"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define a test suite
type GameXMLTestSuite struct {
	suite.Suite
	gameXMLs *gamexml.GameXMLs
}

// SetupTest runs before each test
func (suite *GameXMLTestSuite) SetupTest() {
	suite.gameXMLs = gamexml.New()
}

// TearDownTest runs after each test
func (suite *GameXMLTestSuite) TearDownTest() {
	suite.gameXMLs.Reset()
}

// Test adding a GameXML with a valid ID
func (suite *GameXMLTestSuite) TestAddValidGameXML() {
	g := gamexml.GameXML{Id: "test_id"}
	err := suite.gameXMLs.Add(g)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), g, (*suite.gameXMLs)["test_id"])
}

// Test adding a GameXML with an empty ID
func (suite *GameXMLTestSuite) TestAddInvalidGameXML() {
	g := gamexml.GameXML{Id: ""}
	err := suite.gameXMLs.Add(g)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), gamexml.ErrNoId, err)
}

// Test getting a GameXML with a valid ID
func (suite *GameXMLTestSuite) TestGetValidGameXML() {
	g := gamexml.GameXML{Id: "test_id"}
	_ = suite.gameXMLs.Add(g)
	result, ok := suite.gameXMLs.Get("test_id")
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), g, result)
}

// Test getting a GameXML with an invalid ID
func (suite *GameXMLTestSuite) TestGetInvalidGameXML() {
	result, ok := suite.gameXMLs.Get("invalid_id")
	assert.False(suite.T(), ok)
	assert.Equal(suite.T(), gamexml.GameXML{}, result)
}

// Test resetting the GameXMLs map
func (suite *GameXMLTestSuite) TestResetGameXMLs() {
	g := gamexml.GameXML{Id: "test_id"}
	_ = suite.gameXMLs.Add(g)
	suite.gameXMLs.Reset()
	_, ok := suite.gameXMLs.Get("test_id")
	assert.False(suite.T(), ok)
}

// Test adding multiple GameXMLs
func (suite *GameXMLTestSuite) TestAddMultipleGameXMLs() {
	g1 := gamexml.GameXML{Id: "test_id1"}
	g2 := gamexml.GameXML{Id: "test_id2"}
	err1 := suite.gameXMLs.Add(g1)
	err2 := suite.gameXMLs.Add(g2)
	assert.Nil(suite.T(), err1)
	assert.Nil(suite.T(), err2)
	assert.Equal(suite.T(), g1, (*suite.gameXMLs)["test_id1"])
	assert.Equal(suite.T(), g2, (*suite.gameXMLs)["test_id2"])
}

// Test overwriting an existing GameXML
func (suite *GameXMLTestSuite) TestOverwriteGameXML() {
	g1 := gamexml.GameXML{Id: "test_id"}
	g2 := gamexml.GameXML{Id: "test_id"}
	_ = suite.gameXMLs.Add(g1)
	err := suite.gameXMLs.Add(g2)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), g2, (*suite.gameXMLs)["test_id"])
}

// In order to run the test suite
func TestGameXMLTestSuite(t *testing.T) {
	suite.Run(t, new(GameXMLTestSuite))
}

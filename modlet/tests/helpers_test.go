package modlet_test

import (
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	assert.True(t, modlet.HasPrefix("terrStone", "terr"))
	assert.True(t, modlet.HasPrefix("terrDirt", "terr"))
	assert.False(t, modlet.HasPrefix("rockStone", "terr"))
	assert.False(t, modlet.HasPrefix("", "terr"))
}

func TestHasSuffix(t *testing.T) {
	assert.True(t, modlet.HasSuffix("terrStone", "Stone"))
	assert.True(t, modlet.HasSuffix("rockStone", "Stone"))
	assert.False(t, modlet.HasSuffix("terrDirt", "Stone"))
}

func TestMatch(t *testing.T) {
	assert.True(t, modlet.Match("terrStone", `^terr`))
	assert.True(t, modlet.Match("plantedCorn", `^planted`))
	assert.True(t, modlet.Match("blockShapes", `Shapes$`))
	assert.False(t, modlet.Match("terrStone", `^planted`))
}

func TestNotMatch(t *testing.T) {
	assert.True(t, modlet.NotMatch("terrStone", `^planted`))
	assert.False(t, modlet.NotMatch("plantedCorn", `^planted`))
}

func TestMultValue(t *testing.T) {
	assert.Equal(t, "83", modlet.MultValue("55", 1.5)) // 55 * 1.5 = 82.5 → rounds to 83
	assert.Equal(t, "8", modlet.MultValue("5", 1.5))
	assert.Equal(t, "15,30,45", modlet.MultValue("10,20,30", 1.5))
	assert.Equal(t, "0", modlet.MultValue("", 1.5))
}

func TestProbMult(t *testing.T) {
	assert.Equal(t, "0.45", modlet.ProbMult("0.3", 1.5))
	assert.Equal(t, "1", modlet.ProbMult("0.8", 1.5))
	assert.Equal(t, "1", modlet.ProbMult("0.9", 2.0))
	assert.Equal(t, "1", modlet.ProbMult("1.0", 1.5))
}

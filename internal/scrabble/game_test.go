package scrabble

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/move"
	"math/rand"
	"testing"
)

func TestGame_AICopy(t *testing.T) {
	rand.Seed(0)
	d := dict.NewNode()
	d.Insert("DO")
	p1 := NewComputerPlayer("P1", NewLongestStrategy)
	p2 := NewComputerPlayer("P2", NewLongestStrategy)
	g := NewGame(d, p1, p2)
	assert.Equal(t, "P1", g.CurrentPlayer().Name())
	assert.Equal(t, 0, g.Round)
	g2 := g.AICopy(NewLongestStrategy)
	assert.Equal(t, "P1", g2.CurrentPlayer().Name())
	assert.Equal(t, 0, g2.Round)
	_, err := g2.playMove(move.Move{
		Row:  7,
		Col:  7,
		Dir:  0,
		Word: "DO",
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, g.Round)
	assert.Equal(t, 1, g2.Round)
}

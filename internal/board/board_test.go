package board

import (
	"github.com/stretchr/testify/assert"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/move"
	"github.com/tmazeika/scrabble-go/internal/rules"
	"testing"
)

func abcBoard() *Board {
	b := New(3)
	for i := 0; i < 3*3; i++ {
		b.AtIdx(i).Set(Letter('A' + i))
	}
	return b
}

func TestBoard_At(t *testing.T) {
	b := abcBoard()
	letterIdx := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			assert.Equal(t, Letter('A'+letterIdx), b.At(i, j).Letter())
			letterIdx++
		}
	}
	letterIdx = 0
	b = b.Transposed()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			assert.Equal(t, Letter('A'+letterIdx), b.At(j, i).Letter())
			letterIdx++
		}
	}
}

func TestBoard_Points(t *testing.T) {
	b := New(rules.BoardSize)
	word1 := Word("HORN")
	word2 := Word("FARM")
	word3 := Word("PASTE")
	word4 := Word("MOB")
	word5 := Word("BIT")
	m1 := move.Move{Row: 7, Col: 5, Dir: move.DirAcross, Word: word1}
	m2 := move.Move{Row: 5, Col: 7, Dir: move.DirDown, Word: word2}
	m3 := move.Move{Row: 9, Col: 5, Dir: move.DirAcross, Word: word3}
	m4 := move.Move{Row: 8, Col: 7, Dir: move.DirAcross, Word: word4}
	m5 := move.Move{Row: 10, Col: 4, Dir: move.DirAcross, Word: word5}

	assert.Equal(t, 14, b.Points(m1))
	b.SetAcross(m1.Row, m1.Col, m1.Word)

	bt := b.Transposed()
	m2t := m2.Transposed()
	assert.Equal(t, 9, bt.Points(m2t))
	bt.SetAcross(m2t.Row, m2t.Col, m2t.Word)

	assert.Equal(t, 25, b.Points(m3))
	b.SetAcross(m3.Row, m3.Col, m3.Word)

	assert.Equal(t, 16, b.Points(m4))
	b.SetAcross(m4.Row, m4.Col, m4.Word)

	assert.Equal(t, 16, b.Points(m5))
	b.SetAcross(m5.Row, m5.Col, m5.Word)
}

package board

import (
	"github.com/stretchr/testify/assert"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	"testing"
)

func TestTile_Up(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Letter('B'), b.At(1, 1).Up().Letter())
	b = b.Transposed()
	assert.Equal(t, Letter('G'), b.At(1, 2).Up().Letter())
}

func TestTile_UpN(t *testing.T) {
	b := abcBoard()
	assert.Nil(t, b.At(1, 1).UpN(2))
	assert.Equal(t, Letter('B'), b.At(1, 1).UpN(1).Letter())
	b = b.Transposed()
	assert.Nil(t, b.At(1, 1).UpN(2))
	assert.Equal(t, Letter('D'), b.At(1, 1).UpN(1).Letter())
}

func TestTile_Down(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Letter('H'), b.At(1, 1).Down().Letter())
	b = b.Transposed()
	assert.Equal(t, Letter('H'), b.At(0, 2).Down().Letter())
}

func TestTile_DownN(t *testing.T) {
	b := abcBoard()
	assert.Nil(t, b.At(1, 1).DownN(2))
	assert.Equal(t, Letter('H'), b.At(1, 1).DownN(1).Letter())
	b = b.Transposed()
	assert.Nil(t, b.At(1, 1).DownN(2))
	assert.Equal(t, Letter('F'), b.At(1, 1).DownN(1).Letter())
}

func TestTile_Left(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Letter('A'), b.At(0, 1).Left().Letter())
	b = b.Transposed()
	assert.Equal(t, Letter('D'), b.At(0, 2).Left().Letter())
}

func TestTile_LeftN(t *testing.T) {
	b := abcBoard()
	assert.Nil(t, b.At(1, 1).LeftN(2))
	assert.Equal(t, Letter('D'), b.At(1, 1).LeftN(1).Letter())
	b = b.Transposed()
	assert.Nil(t, b.At(1, 1).LeftN(2))
	assert.Equal(t, Letter('B'), b.At(1, 1).LeftN(1).Letter())
}

func TestTile_Right(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Letter('C'), b.At(0, 1).Right().Letter())
	b = b.Transposed()
	assert.Equal(t, Letter('F'), b.At(2, 0).Right().Letter())
}

func TestTile_RightN(t *testing.T) {
	b := abcBoard()
	assert.Nil(t, b.At(1, 1).RightN(2))
	assert.Equal(t, Letter('F'), b.At(1, 1).RightN(1).Letter())
	b = b.Transposed()
	assert.Nil(t, b.At(1, 1).RightN(2))
	assert.Equal(t, Letter('H'), b.At(1, 1).RightN(1).Letter())
}

func TestTile_GatherUp(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Word("B"), b.At(1, 1).GatherUp())
	b = b.Transposed()
	assert.Equal(t, Word("DE"), b.At(2, 1).GatherUp())
}

func TestTile_GatherDown(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Word("H"), b.At(1, 1).GatherDown())
	b = b.Transposed()
	assert.Equal(t, Word("EF"), b.At(0, 1).GatherDown())
}

func TestTile_GatherLeft(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Word("D"), b.At(1, 1).GatherLeft())
	b = b.Transposed()
	assert.Equal(t, Word("AD"), b.At(0, 2).GatherLeft())
}

func TestTile_GatherRight(t *testing.T) {
	b := abcBoard()
	assert.Equal(t, Word("F"), b.At(1, 1).GatherRight())
	b = b.Transposed()
	assert.Equal(t, Word("DG"), b.At(0, 0).GatherRight())
}

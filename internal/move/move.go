package move

import (
	"fmt"
	. "github.com/tmazeika/scrabble-go/internal/dict"
)

type Dir int

const (
	DirAcross Dir = iota
	DirDown
)

type Move struct {
	Skip bool
	Row  int
	Col  int
	Dir  Dir
	Word Word
}

func (m Move) Transposed() Move {
	m.Row, m.Col = m.Col, m.Row
	if m.Dir == DirAcross {
		m.Dir = DirDown
	} else {
		m.Dir = DirAcross
	}
	return m
}

func (m Move) String() string {
	if m.Skip {
		return "skip"
	}
	dirStr := "across"
	if m.Dir == DirDown {
		dirStr = "down"
	}
	return fmt.Sprintf("(%x,%x) %s: %s", m.Row, m.Col, dirStr, m.Word)
}

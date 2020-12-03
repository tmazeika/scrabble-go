package board

import (
	. "github.com/tmazeika/scrabble-go/internal/dict"
	"strings"
)

type Tile struct {
	letter      *Letter
	board       *Board
	i           int
	crossCheckX *map[Letter]struct{}
	crossCheckY *map[Letter]struct{}
}

func newTile(letter Letter, board *Board, i int) *Tile {
	var nilCrossCheckX map[Letter]struct{} = nil
	var nilCrossCheckY map[Letter]struct{} = nil
	return &Tile{
		letter:      &letter,
		board:       board,
		i:           i,
		crossCheckX: &nilCrossCheckX,
		crossCheckY: &nilCrossCheckY,
	}
}

func (t *Tile) transposed(board *Board, i int) *Tile {
	return &Tile{
		letter:      t.letter,
		board:       board,
		i:           i,
		crossCheckX: t.crossCheckY,
		crossCheckY: t.crossCheckX,
	}
}

func (t *Tile) Letter() Letter {
	return *t.letter
}

func (t *Tile) Set(letter Letter) {
	*t.letter = letter
	*t.crossCheckX = nil
	*t.crossCheckY = nil
	*t.board.staleX = true
	*t.board.staleY = true
}

func (t *Tile) Row() int {
	return t.i / t.board.size
}

func (t *Tile) Col() int {
	return t.i % t.board.size
}

func (t *Tile) Empty() bool {
	return t == nil || *t.letter == 0
}

func (t *Tile) EmptyAround() bool {
	return t.Up().Empty() && t.Down().Empty() &&
		t.Left().Empty() && t.Right().Empty()
}

func (t *Tile) Center() bool {
	return t.i == len(t.board.tiles)/2
}

func (t *Tile) XAnchor() bool {
	return t.Empty() && (!t.Left().Empty() || !t.Right().Empty())
}

func (t *Tile) YAnchor() bool {
	return t.Empty() && (!t.Up().Empty() || !t.Down().Empty())
}

func (t *Tile) UpN(n int) *Tile {
	if t.i/t.board.size < n {
		return nil
	}
	return t.board.tiles[t.i-n*t.board.size]
}

func (t *Tile) Up() *Tile {
	return t.UpN(1)
}

func (t *Tile) DownN(n int) *Tile {
	if t.i/t.board.size >= t.board.size-n {
		return nil
	}
	return t.board.tiles[t.i+n*t.board.size]
}

func (t *Tile) Down() *Tile {
	return t.DownN(1)
}

func (t *Tile) LeftN(n int) *Tile {
	if t.i%t.board.size < n {
		return nil
	}
	return t.board.tiles[t.i-n]
}

func (t *Tile) Left() *Tile {
	return t.LeftN(1)
}

func (t *Tile) RightN(n int) *Tile {
	if t.i%t.board.size >= t.board.size-n {
		return nil
	}
	return t.board.tiles[t.i+n]
}

func (t *Tile) Right() *Tile {
	return t.RightN(1)
}

func (t *Tile) GatherUp() Word {
	return t.gather((*Tile).Up).Reverse()
}

func (t *Tile) GatherDown() Word {
	return t.gather((*Tile).Down)
}

func (t *Tile) GatherLeft() Word {
	return t.gather((*Tile).Left).Reverse()
}

func (t *Tile) GatherRight() Word {
	return t.gather((*Tile).Right)
}

func (t *Tile) gather(nextFn func(*Tile) *Tile) Word {
	var buf strings.Builder
	for cur := nextFn(t); !cur.Empty(); cur = nextFn(cur) {
		buf.WriteRune(rune(cur.Letter()))
	}
	return Word(buf.String())
}

func (t *Tile) InYCrossCheck(l Letter) bool {
	if *t.board.staleY {
		panic("stale cross-check")
	}
	if *t.crossCheckY == nil {
		return true
	}
	_, ok := (*t.crossCheckY)[l]
	return ok
}

func (t *Tile) Premium() (factor int, word bool) {
	if t.board.size != 15 {
		panic("premium squares are undefined for boards of non-standard size")
	}
	if !t.Empty() {
		return 1, false
	}
	// The premium squares are split into 4 quadrants of the board, so we only
	// need to define the premium squares for 1 quadrant: normalize row and col
	// to be in that 1 quadrant.
	row, col := t.board.idxToRowCol(t.i)
	if row > 7 {
		row = 14 - row
	}
	if col > 7 {
		col = 14 - col
	}
	switch row {
	case 0:
		switch col {
		case 0:
			return 3, true
		case 3:
			return 2, false
		case 7:
			return 3, true
		}
	case 1:
		switch col {
		case 1:
			return 2, true
		case 5:
			return 3, false
		}
	case 2:
		switch col {
		case 2:
			return 2, true
		case 6:
			return 2, false
		}
	case 3:
		switch col {
		case 0:
			return 2, false
		case 3:
			return 2, true
		case 7:
			return 2, false
		}
	case 4:
		switch col {
		case 4:
			return 2, true
		}
	case 5:
		switch col {
		case 1:
			return 3, false
		case 5:
			return 3, false
		}
	case 6:
		switch col {
		case 2:
			return 2, false
		case 6:
			return 2, false
		}
	case 7:
		switch col {
		case 0:
			return 3, true
		case 3:
			return 2, false
		case 7:
			return 2, true
		}
	}
	return 1, false
}

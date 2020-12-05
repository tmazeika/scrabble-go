package board

import (
	"fmt"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"github.com/tmazeika/scrabble-go/internal/rules"
	"strings"
)

type Board struct {
	staleX *bool
	staleY *bool
	size   int
	tiles  []*Tile
}

func New(size int) *Board {
	if size < 1 {
		panic("nonpositive board size")
	}
	staleX := false
	staleY := false
	b := Board{
		staleX: &staleX,
		staleY: &staleY,
		size:   size,
		tiles:  make([]*Tile, size*size),
	}
	for i := range b.tiles {
		b.tiles[i] = newTile(0, &b, i)
	}
	return &b
}

func (b *Board) Copy() *Board {
	staleY := *b.staleY
	staleX := *b.staleX
	tiles := make([]*Tile, len(b.tiles))
	b2 := Board{
		staleX: &staleX,
		staleY: &staleY,
		size:   b.size,
		tiles:  tiles,
	}
	for i, t := range b.tiles {
		letter := *t.letter
		var crossCheckX map[Letter]struct{}
		var crossCheckY map[Letter]struct{}
		if *t.crossCheckX != nil {
			crossCheckX = make(map[Letter]struct{}, len(*t.crossCheckX))
			for k, v := range *t.crossCheckX {
				crossCheckX[k] = v
			}
		}
		if *t.crossCheckY != nil {
			crossCheckY = make(map[Letter]struct{}, len(*t.crossCheckY))
			for k, v := range *t.crossCheckY {
				crossCheckY[k] = v
			}
		}
		tiles[i] = &Tile{
			letter:      &letter,
			board:       &b2,
			i:           i,
			crossCheckX: &crossCheckX,
			crossCheckY: &crossCheckY,
		}
	}
	return &b2
}

func (b *Board) Transposed() *Board {
	tiles := make([]*Tile, len(b.tiles))
	b2 := Board{
		staleX: b.staleY,
		staleY: b.staleX,
		size:   b.size,
		tiles:  tiles,
	}
	for i, t := range b.tiles {
		i2 := (i%b.size)*b.size + i/b.size
		tiles[i2] = t.transposed(&b2, i2)
	}
	return &b2
}

func (b *Board) SetYCrossChecks(dict *Node) {
	if !*b.staleY {
		return
	}
	for _, t := range b.tiles {
		if !t.YAnchor() {
			*t.crossCheckY = nil
			continue
		}
		m := make(map[Letter]struct{})
		*t.crossCheckY = m
		above := t.GatherUp()
		below := t.GatherDown()
		for l, n := range dict.Search(above).Edges() {
			if n.Search(below).Accept() {
				m[l] = struct{}{}
			}
		}
	}
	*b.staleY = false
}

func (b *Board) FitsAcross(row, col, n int) bool {
	return row >= 0 && row < b.size && col >= 0 && col < b.size &&
		col+n <= b.size
}

func (b *Board) SetAcross(row, col int, word Word) {
	if !b.FitsAcross(row, col, len(word)) {
		panic("out of bounds")
	}
	i := 0
	for cur := b.At(row, col); i < len(word); cur = cur.Right() {
		cur.Set(Letter(word[i]))
		i++
	}
}

func (b *Board) Points(m Move) int {
	if m.Dir == DirDown {
		b = b.Transposed()
		m = m.Transposed()
	}
	if m.Skip || len(m.Word) == 0 {
		return 0
	}

	var points int
	sumPoints := func(words ...Word) int {
		var sum int
		for _, w := range words {
			for _, l := range w {
				sum += Letter(l).Points()
			}
		}
		return sum
	}

	// Sum across word.
	t := b.At(m.Row, m.Col)
	left := t.GatherLeft()
	right := t.RightN(len(m.Word) - 1).GatherRight()
	// If the across word size is == 1, then it must simply be an extension to
	// a down word. There are no words of length == 1, so don't count them.
	if len(left)+len(m.Word)+len(right) > 1 {
		wordFactor := 1
		var midPoints int
		for i, l := range m.Word {
			factor, word := t.RightN(i).Premium()
			lp := Letter(l).Points()
			if word {
				wordFactor *= factor
				midPoints += lp
			} else {
				midPoints += factor * lp
			}
		}
		points += wordFactor * (midPoints + sumPoints(left, right))
	}

	// Sum down word(s).
	for i, l := range m.Word {
		t := t.RightN(i)
		if !t.Empty() {
			continue
		}
		above := t.GatherUp()
		below := t.GatherDown()
		if len(above)+len(below) == 0 {
			continue
		}
		wordFactor := 1
		midPoints := Letter(l).Points()
		if factor, word := t.Premium(); word {
			wordFactor = factor
		} else {
			midPoints *= factor
		}
		points += wordFactor * (midPoints + sumPoints(above, below))
	}
	if len(m.Word) == rules.RackSize {
		points += rules.BingoPremium
	}
	return points
}

func (b *Board) Anchors() []*Tile {
	var tiles []*Tile
	for _, t := range b.tiles {
		if t.XAnchor() || t.YAnchor() {
			tiles = append(tiles, t)
		}
	}
	return tiles
}

func (b *Board) At(row, col int) *Tile {
	return b.tiles[b.rowColToIdx(row, col)]
}

func (b *Board) AtIdx(i int) *Tile {
	// For side-effects only:
	b.idxToRowCol(i)
	return b.tiles[i]
}

func (b *Board) Center() *Tile {
	return b.tiles[len(b.tiles)/2]
}

func (b *Board) rowColToIdx(row, col int) int {
	if row < 0 || row >= b.size || col < 0 || col >= b.size {
		panic("out of bounds")
	}
	return row*b.size + col
}

func (b *Board) idxToRowCol(i int) (row, col int) {
	if i < 0 || i >= len(b.tiles) {
		panic("out of bounds")
	}
	return i / b.size, i % b.size
}

func (b *Board) String() string {
	var buf strings.Builder
	buf.WriteString("  ")
	for i := 0; i < b.size; i++ {
		buf.WriteString(fmt.Sprintf("%x", i))
		if i < b.size-1 {
			buf.WriteRune(' ')
		} else {
			buf.WriteRune('\n')
		}
	}
	for i, t := range b.tiles {
		if i%b.size == 0 {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(fmt.Sprintf("%x ", i/b.size))
		} else {
			buf.WriteRune(' ')
		}
		if t.Empty() {
			buf.WriteRune('-')
		} else {
			buf.WriteRune(rune(*t.letter))
		}
	}
	return buf.String()
}

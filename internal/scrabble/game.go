package scrabble

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/bag"
	"github.com/tmazeika/scrabble-go/internal/board"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/move"
	"github.com/tmazeika/scrabble-go/internal/rules"
	"math/rand"
	"strings"
)

type Game struct {
	Bag     *bag.Bag
	Board   *board.Board
	Players []Player
	Dict    *Node
	Round   int
}

func NewGame(rnd *rand.Rand, dict *Node, players ...Player) *Game {
	if len(players) < 0 {
		panic("nonpositive player count")
	}
	b := bag.New(rnd)
	for _, p := range players {
		p.DrawFrom(b)
	}
	return &Game{
		Bag:     b,
		Board:   board.New(rules.BoardSize),
		Players: players,
		Dict:    dict,
	}
}

func (g *Game) PlayRound() error {
	player := g.CurrentPlayer()
	m := player.Play(g)
	fmt.Println()
	if m.Skip {
		fmt.Printf("Skipping %s's turn.\n", player.Name())
		g.Round++
		return nil
	}
	b := g.Board
	// Normalize.
	if m.Dir == move.DirDown {
		b = b.Transposed()
		m = m.Transposed()
	}
	b.SetYCrossChecks(g.Dict)

	// Validation.
	if !b.FitsAcross(m.Row, m.Col, len(m.Word)) {
		return fmt.Errorf("move would fall off the board: %v", m)
	}
	if !makesValidWords(g.Dict, b, m) {
		return fmt.Errorf("invalid word(s) would be created: %v", m)
	}
	needed := neededFromRack(b, m)
	if len(needed) == 0 {
		return fmt.Errorf("must put down at least one letter from the rack: %v", m)
	}
	if !player.InRack(needed) {
		return fmt.Errorf("required letters %q are not in rack: %v", needed, m)
	}
	if g.Round == 0 && !passesThroughCenter(b, m) {
		return fmt.Errorf("first move must pass through the center: %v", m)
	}
	if g.Round > 0 && !touchesAnything(b, m) {
		return fmt.Errorf("move must build off an existing move: %v", m)
	}

	// Perform.
	points := b.Points(m)
	fmt.Printf("You scored %d points!\n", points)
	player.AddPoints(points)
	player.UseRack(needed)
	player.DrawFrom(g.Bag)
	b.SetAcross(m.Row, m.Col, m.Word)
	g.Round++
	return nil
}

func (g *Game) Over() bool {
	if !g.Bag.Empty() {
		return false
	}
	for _, p := range g.Players {
		if len(p.Rack()) == 0 {
			return true
		}
	}
	for _, p := range g.Players {
		if len(AllMoves(g.Dict, g.Board, p.Rack())) > 0 {
			return false
		}
	}
	return true
}

func (g *Game) String() string {
	var buf strings.Builder
	buf.WriteString(g.Board.String())
	for _, player := range g.Players {
		buf.WriteRune('\n')
		buf.WriteString(player.String())
	}
	return buf.String()
}

func (g *Game) CurrentPlayer() Player {
	return g.Players[g.Round%len(g.Players)]
}

func makesValidWords(dict *Node, b *board.Board, m move.Move) bool {
	t := b.At(m.Row, m.Col)
	// Valid down.
	for i, l := range m.Word {
		l := Letter(l)
		t := t.RightN(i)
		if !t.Empty() && t.Letter() != l || !t.InYCrossCheck(l) {
			return false
		}
	}
	// Valid across.
	left := t.GatherLeft()
	right := t.RightN(len(m.Word) - 1).GatherRight()
	return dict.Search(left + m.Word + right).Accept()
}

func neededFromRack(b *board.Board, m move.Move) []Letter {
	t := b.At(m.Row, m.Col)
	needed := make([]Letter, 0, len(m.Word))
	for i, l := range m.Word {
		if t.RightN(i).Empty() {
			needed = append(needed, Letter(l))
		}
	}
	return needed
}

func touchesAnything(b *board.Board, m move.Move) bool {
	t := b.At(m.Row, m.Col)
	for i := range m.Word {
		if !t.RightN(i).EmptyAround() {
			return true
		}
	}
	return false
}

func passesThroughCenter(b *board.Board, m move.Move) bool {
	t := b.At(m.Row, m.Col)
	for i := range m.Word {
		if t.RightN(i).Center() {
			return true
		}
	}
	return false
}

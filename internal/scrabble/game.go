package scrabble

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/bag"
	"github.com/tmazeika/scrabble-go/internal/board"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"github.com/tmazeika/scrabble-go/internal/rules"
	"strings"
)

type Game struct {
	Bag     *bag.Bag
	Board   *board.Board
	Players []Player
	Dict    *Node
	Round   int

	over bool
}

func NewGame(dict *Node, players ...Player) *Game {
	if len(players) < 0 {
		panic("nonpositive player count")
	}
	b := bag.New()
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

func (g *Game) AICopy(strategy StrategyFunc) *Game {
	players := make([]Player, len(g.Players))
	for i, p := range g.Players {
		players[i] = p.CopyAsAI(strategy)
	}
	return &Game{
		Bag:     g.Bag.Copy(),
		Board:   g.Board.Copy(),
		Players: players,
		Dict:    g.Dict,
		Round:   g.Round,
		over:    g.over,
	}
}

func (g *Game) PlayRound() (string, error) {
	return g.playMove(g.CurrentPlayer().Play(g))
}

func (g *Game) playMove(m Move) (string, error) {
	player := g.CurrentPlayer()
	if m.Skip {
		g.Round++
		return fmt.Sprintf("\nSkipping %s's turn.\n", player.Name()), nil
	}
	b := g.Board
	// Normalize.
	if m.Dir == DirDown {
		b = b.Transposed()
		m = m.Transposed()
	}
	b.SetYCrossChecks(g.Dict)

	// Validation.
	if !b.FitsAcross(m.Row, m.Col, len(m.Word)) {
		return "", fmt.Errorf("move would fall off the board: %v", m)
	}
	if !makesValidWords(g.Dict, b, m) {
		return "", fmt.Errorf("invalid word(s) would be created: %v", m)
	}
	needed := neededFromRack(b, m)
	if len(needed) == 0 {
		return "", fmt.Errorf("must put down at least one letter from the rack: %v", m)
	}
	if !player.InRack(needed) {
		return "", fmt.Errorf("required letters %q are not in rack: %v", needed, m)
	}
	if g.Board.Center().Empty() && !passesThroughCenter(b, m) {
		return "", fmt.Errorf("first move must pass through the center: %v", m)
	}
	if !g.Board.Center().Empty() && !touchesAnything(b, m) {
		return "", fmt.Errorf("move must build off an existing move: %v", m)
	}

	// Perform.
	points := b.Points(m)
	player.AddPoints(points)
	player.UseRack(needed)
	player.DrawFrom(g.Bag)
	b.SetAcross(m.Row, m.Col, m.Word)
	g.Round++
	return fmt.Sprintf("\n%s scored %d points!\n", player.Name(), points), nil
}

func (g *Game) Over() bool {
	defer g.onOver(g.over)
	if !g.Bag.Empty() {
		return false
	}
	for _, p := range g.Players {
		if len(p.Rack()) == 0 {
			g.over = true
			return true
		}
	}
	for _, p := range g.Players {
		if len(AllMoves(g.Dict, g.Board, p.Rack())) > 0 {
			return false
		}
	}
	g.over = true
	return true
}

func (g *Game) onOver(before bool) {
	if before == g.over {
		return
	}
	totalSum := 0
	for _, p := range g.Players {
		rackSum := 0
		for _, l := range p.Rack() {
			rackSum += l.Points()
		}
		p.AddPoints(-rackSum)
		totalSum += rackSum
	}
	for _, p := range g.Players {
		if len(p.Rack()) == 0 {
			p.AddPoints(totalSum)
		}
	}
}

func (g *Game) Winners() []Player {
	if !g.over {
		panic("game may not be over")
	}
	best, bestPoints := []Player{g.Players[0]}, g.Players[0].Points()
	for _, p := range g.Players[1:] {
		points := p.Points()
		if points > bestPoints {
			best, bestPoints = []Player{p}, points
		} else if points == bestPoints {
			best = append(best, p)
		}
	}
	return best
}

func (g *Game) WonBy(name string) int {
	playerPoints := 0
	best := 0
	for _, p := range g.Players {
		points := p.Points()
		if p.Name() == name {
			playerPoints = points
		} else if points > best {
			best = points
		}
	}
	return playerPoints - best
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

func makesValidWords(dict *Node, b *board.Board, m Move) bool {
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

func neededFromRack(b *board.Board, m Move) []Letter {
	t := b.At(m.Row, m.Col)
	needed := make([]Letter, 0, len(m.Word))
	for i, l := range m.Word {
		if t.RightN(i).Empty() {
			needed = append(needed, Letter(l))
		}
	}
	return needed
}

func touchesAnything(b *board.Board, m Move) bool {
	t := b.At(m.Row, m.Col)
	for i := range m.Word {
		if !t.RightN(i).EmptyAround() {
			return true
		}
	}
	return false
}

func passesThroughCenter(b *board.Board, m Move) bool {
	t := b.At(m.Row, m.Col)
	for i := range m.Word {
		if t.RightN(i).Center() {
			return true
		}
	}
	return false
}

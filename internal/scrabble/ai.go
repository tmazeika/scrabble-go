package scrabble

import (
	"github.com/tmazeika/scrabble-go/internal/board"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"math/rand"
	"sync"
)

type StrategyFunc func(game *Game, moves []Move) Move

func RandomStrategy(_ *Game, moves []Move) Move {
	if len(moves) == 0 {
		return Move{Skip: true}
	}
	return moves[rand.Intn(len(moves))]
}

func LongestStrategy(_ *Game, moves []Move) Move {
	if len(moves) == 0 {
		return Move{Skip: true}
	}
	longest := moves[0]
	for _, m := range moves[1:] {
		if len(m.Word) > len(longest.Word) {
			longest = m
		}
	}
	return longest
}

func MostPointsStrategy(game *Game, moves []Move) Move {
	if len(moves) == 0 {
		return Move{Skip: true}
	}
	best, bestPoints := moves[0], game.Board.Points(moves[0])
	for _, m := range moves[1:] {
		points := game.Board.Points(m)
		if points > bestPoints {
			best, bestPoints = m, points
		}
	}
	return best
}

func NewMCTSStrategy(iterations, pickTop int, c float64) StrategyFunc {
	return func(game *Game, moves []Move) Move {
		if len(moves) == 0 {
			return Move{Skip: true}
		}
		return mcts(game, moves, iterations, pickTop, c)
	}
}

type ComputerPlayer struct {
	basePlayer
	strategy StrategyFunc
}

func NewComputerPlayer(name string, strategy StrategyFunc) *ComputerPlayer {
	return &ComputerPlayer{
		basePlayer{
			name: name,
		},
		strategy,
	}
}

func (p *ComputerPlayer) Play(game *Game) Move {
	return p.strategy(game, AllMoves(game.Dict, game.Board, p.rack))
}

func mergeMoves(cs ...<-chan Move) <-chan Move {
	out := make(chan Move)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, in := range cs {
		in := in
		go func() {
			for m := range in {
				out <- m
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func transposeMoves(in <-chan Move) <-chan Move {
	out := make(chan Move)
	go func() {
		for m := range in {
			out <- m.Transposed()
		}
		close(out)
	}()
	return out
}

func AllMoves(dict *Node, b *board.Board, rack []Letter) []Move {
	b.SetYCrossChecks(dict)
	bt := b.Transposed()
	bt.SetYCrossChecks(dict)
	across := getAcrossMoves(dict, b, rack)
	down := transposeMoves(getAcrossMoves(dict, bt, rack))
	var moves []Move
	for m := range mergeMoves(across, down) {
		moves = append(moves, m)
	}
	return moves
}

func getAcrossMoves(dict *Node, b *board.Board, rack []Letter) <-chan Move {
	out := make(chan Move)
	go func() {
		anchors := b.Anchors()
		if len(anchors) == 0 {
			anchors = []*board.Tile{b.Center()}
		}
		var wg sync.WaitGroup
		wg.Add(len(anchors))
		for _, a := range anchors {
			a := a
			go func() {
				defer wg.Done()
				getAcrossAnchorMoves(dict, b, rack, a, getK(a), out)
			}()
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func getK(anchor *board.Tile) int {
	var k int
	for t := anchor.Left(); t != nil && t.Empty() && t.EmptyAround();
	t = t.Left() {
		k++
	}
	return k
}

func getAcrossAnchorMoves(dict *Node, b *board.Board, rack []Letter,
	anchor *board.Tile, k int, out chan<- Move) {
	if anchor.Left().Empty() {
		leftPart(b, anchor, rack, "", dict, k, out)
	} else {
		left := anchor.GatherLeft()
		extendRight(b, anchor, rack, left, dict.Search(left), anchor, out)
	}
}

func leftPart(b *board.Board, anchor *board.Tile, rack []Letter,
	partialWord Word, node *Node, limit int, out chan<- Move) {
	extendRight(b, anchor, rack, partialWord, node, anchor, out)
	if limit > 0 {
		for l, n := range node.Edges() {
			if Contains(rack, l) {
				leftPart(b, anchor, Remove(rack, l), partialWord.Append(l), n,
					limit-1, out)
			}
		}
	}
}

func extendRight(b *board.Board, anchor *board.Tile, rack []Letter,
	partialWord Word, node *Node, square *board.Tile, out chan<- Move) {
	if square.Empty() {
		if node.Accept() && anchor.Col() < square.Col() {
			out <- Move{
				Row:  square.Row(),
				Col:  square.Col() - len(partialWord),
				Dir:  DirAcross,
				Word: partialWord,
			}
		}
		for l, n := range node.Edges() {
			if Contains(rack, l) && square.InYCrossCheck(l) &&
				square.Right() != nil {
				extendRight(b, anchor, Remove(rack, l),
					partialWord.Append(l), n, square.Right(), out)
			}
		}
	} else if n, ok := node.Edges()[square.Letter()];
		ok && square.Right() != nil {
		extendRight(b, anchor, rack, partialWord.Append(square.Letter()), n,
			square.Right(), out)
	}
}

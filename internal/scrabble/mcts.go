package scrabble

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/board"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"math"
	"sort"
	"strings"
	"sync"
)

const (
	C       = math.Sqrt2
	PickTop = 10
)

type MCTSNode struct {
	mu *sync.RWMutex
	wg *sync.WaitGroup

	parent   *MCTSNode
	children []*MCTSNode

	m      Move
	state  *Game
	score  int
	visits int
}

func (n *MCTSNode) String() string {
	return n.string(0)
}

func (n *MCTSNode) string(indent int) string {
	str := strings.Repeat("  ", indent) +
		fmt.Sprintf("(%d/%d) -- %s\n", n.score, n.visits, n.m)
	for _, c := range n.children {
		str += c.string(indent + 1)
	}
	return str
}

func (n *MCTSNode) search(playerName string, release <-chan struct{}) {
	n.mu.RLock()
	leaf := n
	for len(leaf.children) > 0 {
		leaf = leaf.selectChild()
	}
	if leaf.visits > 0 {
		leaf.expand()
		leaf = leaf.selectChild()
	}
	state := leaf.state.AICopy(MostPointsStrategy)
	n.mu.RUnlock()
	leaf.wg.Add(1)
	go func() {
		defer leaf.wg.Done()
		score := rollout(state, playerName)
		leaf.mu.Lock()
		leaf.backPropagate(score)
		leaf.mu.Unlock()
		<-release
	}()
}

func (n *MCTSNode) selectChild() *MCTSNode {
	if len(n.children) == 0 {
		return n
	}
	var best *MCTSNode
	var bestUCB1 float64
	for _, c := range n.children {
		ucb1 := c.ucb1()
		if math.IsInf(ucb1, 1) {
			return c
		}
		if best == nil || ucb1 > bestUCB1 {
			best, bestUCB1 = c, ucb1
		}
	}
	return best
}

func (n *MCTSNode) ucb1() float64 {
	if n.visits == 0 {
		return math.Inf(1)
	}
	scoref := float64(n.score)
	visitsf := float64(n.visits)
	parentVisitsf := float64(n.parent.visits)
	return scoref/visitsf +
		C*math.Sqrt(math.Log(parentVisitsf)/visitsf)
}

func (n *MCTSNode) expand() {
	n.expandExisting(getTopMoves(n.state.Board,
		AllMoves(n.state.Dict, n.state.Board, n.state.CurrentPlayer().Rack())))
}

func (n *MCTSNode) expandExisting(moves []Move) {
	for _, m := range moves {
		child := MCTSNode{
			mu:     n.mu,
			wg:     n.wg,
			parent: n,
			m:      m,
			state:  n.state.AICopy(nil),
		}
		if _, err := child.state.playMove(m); err != nil {
			panic(err)
		}
		n.children = append(n.children, &child)
	}
}

func rollout(state *Game, playerName string) int {
	for !state.Over() {
		if _, err := state.PlayRound(); err != nil {
			panic(err)
		}
	}
	if wonBy := state.WonBy(playerName); wonBy < 0 {
		return -1
	} else if wonBy == 0 {
		return 0
	}
	return 1
}

func (n *MCTSNode) backPropagate(score int) {
	for n2 := n; n2 != nil; n2 = n2.parent {
		n2.score += score
		n2.visits++
	}
}

func (n *MCTSNode) bestChild() *MCTSNode {
	var best *MCTSNode
	var bestVisits int
	for _, c := range n.children {
		if best == nil || c.visits > bestVisits {
			best, bestVisits = c, c.visits
		}
	}
	return best
}

func mcts(state *Game, moves []Move, runtime int) Move {
	playerName := state.CurrentPlayer().Name()
	root := MCTSNode{
		mu:    &sync.RWMutex{},
		wg:    &sync.WaitGroup{},
		state: state.AICopy(nil),
	}
	root.expandExisting(getTopMoves(state.Board, moves))
	sem := make(chan struct{}, 2)
	i := 0
	for ; i < runtime; i++ {
		sem <- struct{}{}
		root.search(playerName, sem)
	}
	root.wg.Wait()
	fmt.Printf("MCTS stats: %d rollouts completed\n", i)
	return root.bestChild().m
}

func getTopMoves(b *board.Board, m []Move) []Move {
	sort.Slice(m, func(i, j int) bool {
		return b.Points(m[i]) > b.Points(m[j])
	})
	return m[:min(PickTop, len(m))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

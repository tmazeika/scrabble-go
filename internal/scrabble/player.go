package scrabble

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/bag"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"github.com/tmazeika/scrabble-go/internal/rules"
	"strings"
)

type Player interface {
	fmt.Stringer
	Name() string
	Points() int
	AddPoints(points int)
	Rack() []Letter
	InRack(letters []Letter) bool
	UseRack(letters []Letter)
	DrawFrom(bag *bag.Bag)
	Play(game *Game) Move
	CopyAsAI(strategy StrategyFunc) Player
}

type basePlayer struct {
	name   string
	points int
	rack   []Letter
}

func (p *basePlayer) Name() string {
	return p.name
}

func (p *basePlayer) Points() int {
	return p.points
}

func (p *basePlayer) AddPoints(points int) {
	p.points += points
}

func (p *basePlayer) Rack() []Letter {
	rack := make([]Letter, len(p.rack))
	copy(rack, p.rack)
	return rack
}

func (p *basePlayer) InRack(letters []Letter) bool {
	rack := p.Rack()
	for _, l := range letters {
		if !Contains(rack, l) {
			return false
		}
		rack = Remove(rack, l)
	}
	return true
}

func (p *basePlayer) UseRack(letters []Letter) {
	for _, l := range letters {
		p.rack = Remove(p.rack, l)
	}
}

func (p *basePlayer) DrawFrom(bag *bag.Bag) {
	p.rack = append(p.rack, bag.Draw(rules.RackSize-len(p.rack))...)
}

func (p *basePlayer) CopyAsAI(strategy StrategyFunc) Player {
	p2 := ComputerPlayer{
		basePlayer: basePlayer{
			name:   p.name,
			points: p.points,
			rack:   make([]Letter, len(p.rack)),
		},
		strategy: strategy,
	}
	copy(p2.rack, p.rack)
	return &p2
}

func (p *basePlayer) String() string {
	ls := make([]string, len(p.rack))
	for i, l := range p.rack {
		ls[i] = string(l)
	}
	return fmt.Sprintf("%s (%d points): %s",
		p.name, p.points, strings.Join(ls, ","))
}

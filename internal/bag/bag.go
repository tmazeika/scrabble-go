package bag

import (
	. "github.com/tmazeika/scrabble-go/internal/dict"
	"math/rand"
	"strings"
)

type Bag struct {
	letters []Letter
}

func New(rnd *rand.Rand) *Bag {
	ls := LettersDist()
	rnd.Shuffle(len(ls), func(i, j int) {
		ls[i], ls[j] = ls[j], ls[i]
	})
	return &Bag{ls}
}

func (b *Bag) Draw(n int) []Letter {
	ls := make([]Letter, min(n, len(b.letters)))
	b.letters = b.letters[copy(ls, b.letters):]
	return ls
}

func (b *Bag) Empty() bool {
	return len(b.letters) == 0
}

func (b *Bag) String() string {
	ls := make([]string, len(b.letters))
	for i, l := range b.letters {
		ls[i] = string(l)
	}
	return strings.Join(ls, ",")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

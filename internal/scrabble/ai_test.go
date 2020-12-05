package scrabble

import (
	"github.com/tmazeika/scrabble-go/internal/dict"
	"math/rand"
	"testing"
)

func BenchmarkComputerPlayer_Play(b *testing.B) {
	println("Here")
	rand.Seed(1)
	d, err := dict.Load(dict.Dict)
	if err != nil {
		panic(err)
	}
	player1 := NewComputerPlayer("P1", NewRandomStrategy())
	player2 := NewComputerPlayer("P2", NewRandomStrategy())
	game := NewGame(d, player1, player2)
	_, err = game.PlayRound()
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = player2.Play(game)
	}
}

package main

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/scrabble"
	"math/rand"
	"os"
	"time"
)

const (
	Iterations = 50
	Trials     = 300
)

var winners = map[string]int{}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := play(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, wins := range winners {
		fmt.Println(wins, "wins for", name)
	}
}

func play() error {
	root, err := dict.Load(dict.Dict)
	if err != nil {
		return err
	}
	for i := 0; i < Trials; i++ {
		player1 := scrabble.NewComputerPlayer("TJ",
			scrabble.NewRandomStrategy())
		player2 := scrabble.NewComputerPlayer("Justine",
			scrabble.NewRandomStrategy())
		game := scrabble.NewGame(root, player1, player2)
		for !game.Over() {
			// fmt.Println(game.String())
			_, err := game.PlayRound()
			if err != nil {
				panic(err)
				// fmt.Printf("Bad move: %v\n", err)
			} else {
				// fmt.Print(s)
			}
		}
		// fmt.Println(game.String())
		// fmt.Println("Game over!")
		for _, p := range game.Winners() {
			fmt.Println("Winner:", p.Name())
			winners[p.Name()]++
		}
	}
	return nil
}

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
	Iterations = 30
	Trials     = 1
)

var winners = map[string]int{}

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := play(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// for name, wins := range winners {
	// 	fmt.Println(wins, "wins for", name)
	// }
}

func play() error {
	root, err := dict.Load(dict.Dict)
	if err != nil {
		return err
	}
	for i := 0; i < Trials; i++ {
		player1 := scrabble.NewComputerPlayer("MostPoints",
			scrabble.MostPointsStrategy)
		player2 := scrabble.NewComputerPlayer("MCTS-AI",
			scrabble.NewMCTSStrategy(Iterations))
		game := scrabble.NewGame(root, player1, player2)
		for !game.Over() {
			s, err := game.PlayRound()
			fmt.Println(game.String())
			if err != nil {
				fmt.Printf("Bad move: %v\n", err)
			} else {
				fmt.Print(s)
			}
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println(game.String())
		// fmt.Println("Game over!")
		// for _, p := range game.Winners() {
		// 	fmt.Println("Winner:", p.Name())
		// 	winners[p.Name()]++
		// }
	}
	return nil
}

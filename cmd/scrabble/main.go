package main

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/scrabble"
	"math/rand"
	"os"
	"time"
)

const Trials = 1

var winners = map[string]int{}
var gameLeads [][]int

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := play(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, wins := range winners {
		fmt.Printf("%d wins for %s\n", wins, name)
	}
	for _, game := range gameLeads {
		for i, lead := range game {
			fmt.Printf("%d", lead)
			if i < len(game) - 1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
}

func play() error {
	root, err := dict.Load(dict.Dict)
	if err != nil {
		return err
	}
	for i := 0; i < Trials; i++ {
		// var leads []int
		player1 := scrabble.NewComputerPlayer("MostPoints",
			scrabble.MostPointsStrategy)
		player2 := scrabble.NewComputerPlayer("MCTS-AI",
			scrabble.NewMCTSStrategy(25, 10, 1.4))
		game := scrabble.NewGame(root, player1, player2)
		for !game.Over() {
			fmt.Println(game.String())
			s, err := game.PlayRound()
			if err != nil {
				fmt.Printf("Bad move: %v\n", err)
			} else {
				fmt.Print(s)
			}
			// if game.CurrentPlayer().Name() == "MCTS-AI-1.5" {
			// 	leads = append(leads, player2.Points() - player1.Points())
			// }
		}
		// gameLeads = append(gameLeads, leads)
		fmt.Println(game.String())
		fmt.Println("Game over!")
		for _, p := range game.Winners() {
			fmt.Println("Winner:", p.Name())
			winners[p.Name()]++
		}
	}
	return nil
}

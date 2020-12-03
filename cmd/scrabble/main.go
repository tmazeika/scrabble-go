package main

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/scrabble"
	"math/rand"
	"os"
)

func main() {
	if err := play(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func play() error {
	rnd := rand.New(rand.NewSource(0))
	root, err := dict.Load(dict.Dict)
	if err != nil {
		return err
	}
	player1 := scrabble.NewComputerPlayer("TJ",
		scrabble.NewLongestStrategy())
	player2 := scrabble.NewComputerPlayer("Justine",
		scrabble.NewMostPointsStrategy())
	game := scrabble.NewGame(rnd, root, player1, player2)
	for !game.Over() {
		fmt.Println(game.String())
		if err := game.PlayRound(); err != nil {
			fmt.Printf("Bad move: %v\n", err)
			continue
		}
	}
	fmt.Println(game.String())
	fmt.Println()
	fmt.Println("Game over!")
	return nil
}

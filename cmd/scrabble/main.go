package main

import (
	"fmt"
	"github.com/tmazeika/scrabble-go/internal/dict"
	"github.com/tmazeika/scrabble-go/internal/scrabble"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	if err := play(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func play() error {
	root, err := dict.Load(dict.Dict)
	if err != nil {
		return err
	}
	player1 := scrabble.NewComputerPlayer("TJ",
		scrabble.NewMCTSStrategy(15 * time.Second))
	player2 := scrabble.NewComputerPlayer("Justine",
		scrabble.MostPointsStrategy)
	game := scrabble.NewGame(root, player1, player2)
	for !game.Over() {
		fmt.Println(game.String())
		s, err := game.PlayRound()
		if err != nil {
			fmt.Printf("Bad move: %v\n", err)
		} else {
			fmt.Print(s)
		}
	}
	fmt.Println(game.String())
	fmt.Println("Game over!")
	return nil
}

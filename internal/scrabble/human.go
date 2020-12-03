package scrabble

import (
	"bufio"
	"fmt"
	. "github.com/tmazeika/scrabble-go/internal/dict"
	. "github.com/tmazeika/scrabble-go/internal/move"
	"os"
	"strconv"
	"strings"
)

type HumanPlayer struct {
	basePlayer
}

func NewHumanPlayer(name string) *HumanPlayer {
	return &HumanPlayer{
		basePlayer{
			name: name,
		},
	}
}

func (p *HumanPlayer) Play(*Game) Move {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Move for %s [s(kip)|(row,col,a(cross)|d(own),letters...)]: ",
			p.name)
		var moveStr string
		if scanner.Scan() {
			moveStr = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if moveStr == "s" || moveStr == "skip" {
			return Move{Skip: true}
		}
		split := strings.SplitN(moveStr, ",", 4)
		if len(split) != 4 {
			fmt.Println("Bad move format. Try again...")
			continue
		}
		row, err := strconv.ParseInt(split[0], 16, 32)
		if err != nil {
			fmt.Println("Invalid row. Try again...")
			continue
		}
		col, err := strconv.ParseInt(split[1], 16, 32)
		if err != nil {
			fmt.Println("Invalid column. Try again...")
			continue
		}
		var dir Dir
		if split[2] == "a" || split[2] == "across" {
			dir = DirAcross
		} else if split[2] == "d" || split[2] == "down" {
			dir = DirDown
		} else {
			fmt.Println("Invalid direction. Try again...")
			continue
		}
		word := Word(split[3])
		if !IsWord(split[3]) {
			fmt.Println("Invalid letters. Try again...")
		}
		return Move{
			Row:  int(row),
			Col:  int(col),
			Dir:  dir,
			Word: word,
		}
	}
}

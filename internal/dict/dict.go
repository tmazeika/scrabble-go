package dict

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const Dict = "dictionary.txt"

type Letter rune

func IsLetter(r rune) bool {
	return 'A' <= r && r <= 'Z' || r == '_'
}

func Contains(letters []Letter, letter Letter) bool {
	for _, l := range letters {
		if l == letter {
			return true
		}
	}
	return false
}

func Remove(letters []Letter, letter Letter) []Letter {
	ls := make([]Letter, len(letters)-1)
	for i, l := range letters {
		if l == letter {
			copy(ls[:i], letters[:i])
			copy(ls[i:], letters[i+1:])
			return ls
		}
	}
	return letters
}

type Word string

func (w Word) Head() Letter {
	return Letter(string(w)[0])
}

func (w Word) Tail() Word {
	return Word(string(w)[1:])
}

func (w Word) Reverse() Word {
	rs := []rune(w)
	n := len(w)
	for i := 0; i < n/2; i++ {
		rs[i], rs[n-1-i] = rs[n-1-i], rs[i]
	}
	return Word(rs)
}

func (w Word) Append(l Letter) Word {
	return Word(string(w) + string(l))
}

func IsWord(s string) bool {
	for _, r := range s {
		if !IsLetter(r) {
			return false
		}
	}
	return true
}

func (w Word) Without(l Letter) Word {
	i := strings.IndexRune(string(w), rune(l))
	if i == -1 {
		return w
	}
	return w[:i] + w[i+1:]
}

var letterProps = map[Letter]struct {
	points int
	count  int
}{
	'A': {points: 1, count: 9},
	'B': {points: 3, count: 2},
	'C': {points: 3, count: 2},
	'D': {points: 2, count: 4},
	'E': {points: 1, count: 12},
	'F': {points: 4, count: 2},
	'G': {points: 2, count: 3},
	'H': {points: 4, count: 2},
	'I': {points: 1, count: 9},
	'J': {points: 8, count: 1},
	'K': {points: 5, count: 1},
	'L': {points: 1, count: 4},
	'M': {points: 3, count: 2},
	'N': {points: 1, count: 6},
	'O': {points: 1, count: 8},
	'P': {points: 3, count: 2},
	'Q': {points: 10, count: 1},
	'R': {points: 1, count: 6},
	'S': {points: 1, count: 4},
	'T': {points: 1, count: 6},
	'U': {points: 1, count: 4},
	'V': {points: 4, count: 2},
	'W': {points: 4, count: 2},
	'X': {points: 8, count: 1},
	'Y': {points: 4, count: 2},
	'Z': {points: 10, count: 1},
	// TODO: '_': {points: 0, count: 2},
	'_': {points: 0, count: 0},
}

func (l Letter) Points() int {
	p, ok := letterProps[l]
	if !ok {
		panic(fmt.Sprintf("unknown letter %q", l))
	}
	return p.points
}

func LettersDist() []Letter {
	var ls []Letter
	for l, p := range letterProps {
		for i := 0; i < p.count; i++ {
			ls = append(ls, l)
		}
	}
	return ls
}

func Load(filename string) (n *Node, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer safeClose(f, &err)
	scanner := bufio.NewScanner(f)
	n = NewNode()
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) > 0 {
			n.Insert(Word(line))
		}
	}
	return n, scanner.Err()
}

func safeClose(closer io.Closer, err *error) {
	if cerr := closer.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

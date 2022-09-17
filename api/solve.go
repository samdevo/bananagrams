package function

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const DICTLEN int = 178691
const MAXLEN int = 5
const MINLEN int = 3

type Cell struct {
	char byte
	r, c int
}

type board [][]byte

const EMPTY = byte(' ')

type EmptySpace struct {
	cell                    Cell
	spaceBefore, spaceAfter int
	isHorizontal            bool
}

type Game struct {
	board      board
	chars      []byte
	dictionary []string
}

func (b board) print() {
	if b == nil {
		return
	}
	rows := len(b)
	cols := len(b[0])
	fmt.Printf("rows: %d, cols: %d \n--------------\n", rows, cols)
	for _, row := range b {
		fmt.Println(string(row))
	}
	fmt.Println("--------------------")
}

func newGame(chars string, dict string) *Game {
	dictionary := getDictionary("dictionary.txt")
	return &Game{chars: []byte(chars), dictionary: dictionary}
}

func (game *Game) solve() (sol board) {
	return game.solveSetLen(MINLEN, MAXLEN)
}

func (game *Game) solveSetLen(minLen, maxLen int) (sol board) {
	// defer close(boardStream)
	solution := make(chan board)
	fail := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if len(game.board) == 0 {
			game.fromStartingGames(ctx, solution, minLen, maxLen)
		} else {
			game.search(ctx, solution, 0, minLen, maxLen)
		}
		fail <- true
	}()
	select {
	case <-fail:
		fmt.Println("main search failed :(")
	case sol = <-solution:
		fmt.Println("solution received!")
	case <-time.After(30 * time.Second):
		fmt.Println("timeout!")
	}
	return
}

type ValidEntry struct {
	cell         Cell
	entry        string
	cellInd      int
	isHorizontal bool
}

// seeds new games with boards of one character
func (g *Game) fromStartingGames(ctx context.Context, done chan board, minLen, maxLen int) {
	for i, char := range g.chars {
		newBoard := make(board, 1)
		newBoard[0] = []byte{char}
		newChars := make([]byte, len(g.chars)-1)
		for j := 0; j < len(g.chars)-1; j++ {
			if j < i {
				newChars[j] = g.chars[j]
			} else {
				newChars[j] = g.chars[j+1]
			}
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(newChars), func(i, j int) { newChars[i], newChars[j] = newChars[j], newChars[i] })
		newGame := &Game{newBoard, newChars, g.dictionary}
		newGame.search(ctx, done, 1, minLen, maxLen)
	}
}

// func worker(jobChan chan *Game, ctx context.Context, done chan board) {
// 	for game := range jobChan {
// 		game.search(ctx, done)
// 	}
// }

// searches the game for a solution
func (g *Game) search(ctx context.Context, done chan board, depth, minLen, maxLen int) {
	ch := make(chan ValidEntry)
	doneSearching := make(chan bool)
	// quit := make(chan bool)
	emptySpaces := g.board.getSpaces()
	validCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go g.findValidWords(validCtx, emptySpaces, ch, minLen, maxLen)
	visited := make(map[string]bool)
	for {
		select {
		case valid, ok := <-ch:
			if !ok {
				return
			}
			if visited[valid.entry] {
				continue
			}
			visited[valid.entry] = true
			// fmt.Println("found valid entry: ", valid.entry)

			newGame := g.addToBoard(valid)
			// fmt.Println("chars left: ", len(newGame.chars))
			// newGame.board.print()

			if len(newGame.chars) == 0 {
				fmt.Println("found a solution!")
				select {
				case done <- newGame.board:
				case <-ctx.Done():
				}
				return
			}
			// select {
			// case boardStream <- newGame.board:
			// case <-ctx.Done():
			// 	return
			// }
			// fmt.Printf("depth: %d\n", depth)
			newGame.search(ctx, done, depth+1, minLen, maxLen)
		case <-doneSearching:
			break
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			fmt.Println("timeout!")
			return
		}
	}
	// fmt.Println("done searching current game")
}

// given a permutation and an empty space, determines if the permutation can create a ValidWord
func processPerm(ctx context.Context, space EmptySpace, wordLen int, str []byte, validChan chan ValidEntry, dict []string) bool {
	word := str[:wordLen]
	// sstr := string(str)
	// fmt.Println(sstr)
	minInd := 0
	maxInd := wordLen
	if space.spaceBefore != -1 && space.spaceBefore < wordLen {
		maxInd = space.spaceBefore
	}
	if space.spaceAfter != -1 && space.spaceAfter < wordLen {
		minInd = wordLen - space.spaceAfter
	}
	for cellInd := minInd; cellInd <= maxInd; cellInd++ {
		// newWord := []byte(strings.Clone(string(word)))
		// newWordStr := string(append(append(newWord[:cellInd], space.cell.char), newWord[cellInd:]...))

		newWordStr := strings.Join([]string{string(word[:cellInd]), string(word[cellInd:])}, string(space.cell.char))
		// newWord := make([]byte, wordLen+1)
		// copy(newWord, append(append(word[:cellInd], space.cell.char), word[cellInd:]...))

		if validWord(newWordStr, dict) {
			select {
			case <-ctx.Done():
				return false
			case validChan <- ValidEntry{space.cell, newWordStr, cellInd, space.isHorizontal}:
				return true
			}

		}
	}
	return true
}

// given a game, pushes valid words that can be added to the board to the validChan channel
func (g *Game) findValidWords(ctx context.Context, emptySpaces []EmptySpace, validChan chan ValidEntry, minLen, maxLen int) {
	defer close(validChan)
	for _, space := range emptySpaces {
		visited := make(map[string]bool)
		startLen := maxLen - 1
		if space.spaceBefore != -1 && space.spaceAfter != -1 {
			if startLen > space.spaceBefore+space.spaceAfter {
				startLen = space.spaceBefore + space.spaceAfter
			}
		}
		if startLen > len(g.chars) {
			startLen = len(g.chars)
		}
		for wordLen := startLen; wordLen >= minLen-1; wordLen-- {
			if len(g.chars)-wordLen != 0 && len(g.chars)-wordLen < minLen-1 {
				continue
			}
			perm(g.chars, func(newStr []byte) {
				if visited[string(newStr)] {
					return
				}
				visited[string(newStr)] = true
				if !processPerm(ctx, space, wordLen, newStr, validChan, g.dictionary) {
					return
				}
			}, 0, wordLen, ctx)

		}
	}
}

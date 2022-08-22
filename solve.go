package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

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
	for _, row := range b {
		fmt.Println(string(row))
	}
}

func newGame(chars string, dict string) *Game {
	dictionary := getDictionary("dictionary.txt")
	return &Game{chars: []byte(chars), dictionary: dictionary}
}

func (game *Game) solve() [][]byte {
	// var board [][]byte
	// firstWordOptions := findValidWords(chars, dictionary)
	// fmt.Printf("%v\n", firstWordOptions)
	var wg sync.WaitGroup
	solution := make(chan board)
	fail := make(chan bool)
	go func() {
		wg.Add(1)
		game.search(&wg, solution)
		wg.Wait()
		fail <- true
	}()
	select {
	case <-fail:
		fmt.Println("main search failed :(")
		return nil
	case sol := <-solution:
		return sol
	}
}

type ValidEntry struct {
	cell         Cell
	entry        string
	cellInd      int
	isHorizontal bool
}

func (g *Game) search(wg *sync.WaitGroup, done chan board) {
	defer wg.Done()
	ch := make(chan ValidEntry)
	emptySpaces := g.board.getSpaces()
	go g.findValidWords(emptySpaces, ch)
	for valid := range ch {
		// fmt.Println(string(str))
		fmt.Println(valid)
		fmt.Println("____")
	}
}

func (g *Game) findValidWords(emptySpaces []EmptySpace, validChan chan ValidEntry) {
	for _, space := range emptySpaces {
		startLen := MAXLEN
		if space.spaceBefore != -1 && space.spaceAfter != -1 {
			if startLen > space.spaceBefore+space.spaceAfter {
				startLen = space.spaceBefore + space.spaceAfter
			}
		}
		for wordLen := startLen; wordLen >= MINLEN; wordLen-- {
			perm([]byte(g.chars), func(str []byte) {
				word := str[:wordLen]
				minInd := 0
				maxInd := wordLen - 1
				if space.spaceBefore != -1 {
					maxInd = space.spaceBefore
				}
				if space.spaceAfter != -1 {
					minInd = wordLen - space.spaceAfter
				}
				for cellInd := minInd; cellInd <= maxInd; cellInd++ {
					newWord := append(append(word[:cellInd], space.cell.char), word[cellInd:]...)
					if validWord(string(newWord), g.dictionary) {
						validChan <- ValidEntry{space.cell, string(newWord), cellInd, space.isHorizontal}
					}
				}
			}, 0, wordLen)
		}
	}
	// if len(chars) < MINLEN {
	// 	return
	// }

	// startLen := MAXLEN
	// if len(chars) < MAXLEN {
	// 	startLen = len(chars)
	// }
	// for wordLen := startLen; wordLen >= MINLEN; wordLen-- {
	// 	perm([]byte(chars), func(str []byte) {
	// 		if validWord(string(str[:wordLen]), dictionary) {
	// 			fmt.Println(string(str[:wordLen]))
	// 			validChan <- str[:wordLen]
	// 		}
	// 	}, 0, wordLen)
	// }
	// close(validChan)
	// return
}

// func (g *Game) getSpaces

// func addToBoard(str []byte, board [][]byte) [][]byte {
// 	if len(board) == 0 {
// 		return append(board, str)
// 	}
// 	for i, row := range board {
// 		if row[0]
// 	}
// }

// searches the board for the beginning of a word

// func (b Game) findSpaces() (spaces []EmptySpace) {
// 	for i, row := range b {
// 		hCounter := -1
// 		for j, cell := range row {
// 			if cell != 0 {
// 				if hCounter == -1 {
// 					hCounter = 0
// 				} else if hCounter == 0 {

// 				}
// 			}
// 		}
// 	}
// 	return
// }

// func (g *Game) loadSpaces() {
// 	for i, row := range g.board {
// 		lastByte := byte('0')
// 		lastCharCell := nil
// 		leftCounter := -1
// 		rightCounter := -1
// 		horizontalWord := false
// 		for j, cell := range row {
// 			if cell != 0 { // cell is a character
// 				if horizontalWord {
// 					continue
// 				}
// 				if lastChar != byte('0') {
// 					horizontalWord = true
// 					continue
// 				}
// 				lastCharCell := Cell{cell, i, j}
// 				if rightCounter != -1 {
// 					g.spaces = append(g.spaces, EmptySpace{lastCharCell, leftCounter, rightCounter - 1})
// 					leftCounter = 0
// 					rightCounter = -1
// 				}
// 			} else { // blank space
// 				if lastByte {
// 					lastByte = false
// 					if hCounter != 0 {
// 						hCounter = 1
// 					}
// 				}
// 			}
// 			lastByte =
// 		}
// 		if !horizontalWord && rightCounter > 0 {
// 			g.spaces = append(g.spaces, EmptySpace{lastChar})
// 		}
// 	}
// }

// returns a slice of EmptySpace by looping through the board
func (b board) getSpaces() []EmptySpace {
	spaces := b.getHorizontalSpaces(false)
	spaces = append(spaces, b.transposed().getHorizontalSpaces(true)...)
	return spaces
}

func (b board) getHorizontalSpaces(isTransposed bool) (spaces []EmptySpace) {
	for i, row := range b {
		for j, cell := range row {
			if isChar(cell) {
				left, adjLeft := spaceLeft(i, j, b)
				right, adjRight := spaceRight(i, j, b)
				if !adjLeft || !adjRight || (left == 0 && right == 0) {
					continue
				}
				var newSpace EmptySpace
				if !isTransposed {
					newSpace = EmptySpace{Cell{cell, i, j}, left, right, true}
				} else {
					newSpace = EmptySpace{Cell{cell, j, i}, left, right, false}
				}
				spaces = append(spaces, newSpace)
			}
		}
	}
	return
}

func isChar(val byte) bool {
	return val != EMPTY
}

// returns a transposed version of the board
func (b board) transposed() board {
	xl := len(b[0])
	yl := len(b)
	result := make(board, xl)
	for i := range result {
		result[i] = make([]byte, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = b[j][i]
		}
	}
	return result
}

// returns the amount of free space to the left of the given cell
func spaceLeft(r, c int, board [][]byte) (int, bool) {
	if c == 0 {
		return -1, true
	}
	curCol := c
	charSpace := 0
	for {
		curCol--
		if curCol < 0 {
			return -1, true
		}
		curByte := board[r][curCol]
		if curByte == EMPTY {
			if !hasVerticalAdjacent(r, curCol, board) {
				charSpace++
			} else {
				return charSpace, true
			}
		} else {
			return charSpace, false
		}
	}
}

// returns the amount of free space to the right of the given cell
func spaceRight(r, c int, board [][]byte) (int, bool) {
	if c == len(board[0])-1 {
		return -1, true
	}
	curCol := c
	charSpace := 0
	for {
		curCol++
		if curCol == len(board[0]) {
			return -1, true
		}
		curByte := board[r][curCol]
		if curByte == EMPTY {
			if !hasVerticalAdjacent(r, curCol, board) {
				charSpace++
			} else {
				return charSpace, true
			}
		} else {
			return charSpace, false
		}
	}
}

// returns true if the cell has an adjacent character (non-empty) above or below
func hasVerticalAdjacent(r, c int, board [][]byte) bool {
	if r == 0 {
		if len(board) == 1 {
			return false
		}
		if board[r+1][c] != EMPTY {
			return true
		}
	} else if r == len(board)-1 {
		if board[r-1][c] != EMPTY {
			return true
		}
	} else {
		if board[r+1][c] != EMPTY || board[r-1][c] != EMPTY {
			return true
		}
	}
	return false
}

func getDictionary(filename string) (dictionary []string) {
	dictionary = make([]string, DICTLEN)
	ind := 0
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary[ind] = scanner.Text()
		ind++
	}
	return
}

func validWord(word string, dictionary []string) bool {
	//	binary search (words sorted)
	upperBound := DICTLEN - 1
	lowerBound := 0
	for {
		ind := (upperBound + lowerBound) / 2
		switch strings.Compare(word, dictionary[ind]) {
		case 0:
			return true
		case -1:
			upperBound = ind - 1
		case 1:
			lowerBound = ind + 1
		}
		if upperBound < lowerBound {
			return false
		}
	}
}

func findPrefixMatches(prefix string, dictionary []string) (words []string) {
	upperBound := DICTLEN - 1
	lowerBound := 0

	for {
		ind := (upperBound + lowerBound) / 2
		if strings.HasPrefix(dictionary[ind], prefix) {
			words = append(words, dictionary[ind])
			origInd := ind
			for {
				ind++
				if ind >= len(dictionary) || !strings.HasPrefix(dictionary[ind], prefix) {
					break
				}
				words = append(words, dictionary[ind])
			}
			ind = origInd
			for {
				ind--
				if ind < 0 || !strings.HasPrefix(dictionary[ind], prefix) {
					break
				}
				words = append([]string{dictionary[ind]}, words...)
			}
			return
		}
		switch strings.Compare(prefix, dictionary[ind]) {
		case -1:
			upperBound = ind - 1
		case 1:
			lowerBound = ind + 1
		}
		if upperBound < lowerBound {
			return
		}
	}
}

//// Perm calls f with each permutation of a.
//func Perm(a []byte, f func([]byte)) {
//	perm(a, f, 0, len(a))
//}

// gets permutations of string a
// starts permuting at start and stops at the index before stop
// if all is true, the entire string is included, regardless of where stop is
// if all is false, the string up to
func getPermutations(a string, start, stop int, all bool) (rt []string) {
	rt = make([]string, permLen(len(a), stop-start))
	i := 0
	strStop := stop
	strStart := start
	if all {
		strStop = len(a)
		strStart = 0
	}
	if stop < 0 {
		stop = len(a) + stop
	}
	perm([]byte(a), func(str []byte) {
		rt[i] = string(str[strStart:strStop])
		i++
	}, start, stop)
	fmt.Printf("num permutations: %d, should be: %d\n", len(rt), permLen(len(a), stop-start))
	return
}

// Permute the values at index i to stop.
func perm(a []byte, f func([]byte), i, stop int) {
	if i > stop-1 {
		f(a)
		return
	}
	perm(a, f, i+1, stop)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1, stop)
		a[i], a[j] = a[j], a[i]
	}
}

func permLen(n, k int) int {
	return factorial(n) / factorial(n-k)
}

func factorial(x int) int {
	if x <= 0 {
		return 1
	}
	return x * factorial(x-1)
}

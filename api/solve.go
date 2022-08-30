package function

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
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
	// defer close(boardStream)
	solution := make(chan board)
	fail := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if len(game.board) == 0 {
			game.fromStartingGames(ctx, solution, MINLEN, MAXLEN)
		} else {
			game.search(ctx, solution, 0, MINLEN, MAXLEN)
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

func (g *Game) addToBoard(entry ValidEntry) *Game {
	newWord := []byte(entry.entry)
	var curBoard board
	var wordStart int
	var wordRow int
	if !entry.isHorizontal {
		curBoard = g.board.transposed()
		wordStart = entry.cell.r - entry.cellInd
		wordRow = entry.cell.c
	} else {
		curBoard = g.board.copy()
		wordStart = entry.cell.c - entry.cellInd
		wordRow = entry.cell.r
	}
	newColsLeft := wordStart * -1
	for i := 0; i < newColsLeft; i++ {
		for i, row := range curBoard {
			curBoard[i] = append([]byte{EMPTY}, row...)
		}
		wordStart++
	}
	newColsRight := wordStart + len(entry.entry) - len(curBoard[wordRow])
	for i := 0; i < newColsRight; i++ {
		for i, row := range curBoard {
			curBoard[i] = append(row, EMPTY)
		}
	}
	curBoard[wordRow] =
		append(append(
			curBoard[wordRow][:wordStart],
			newWord...),
			curBoard[wordRow][wordStart+len(entry.entry):]...,
		)
	var charsLeft []byte
	for _, char := range g.chars {
		found := false
		for j, charUsed := range newWord {
			if j == entry.cellInd {
				continue
			}
			if char == charUsed {
				newWord[j] = EMPTY
				found = true
				break
			}
		}
		if !found {
			charsLeft = append(charsLeft, char)
		}
	}

	var newBoard board
	if !entry.isHorizontal {
		newBoard = curBoard.transposed()
	} else {
		newBoard = curBoard
	}

	return &Game{newBoard, charsLeft, g.dictionary}
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

// given a permutation of indices and the list of characters left, produce the given word
func permToStr(inds []int, chars []byte) []byte {
	newStr := make([]byte, len(inds))
	for i, ind := range inds {
		newStr[i] = chars[ind]
	}
	return newStr
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
			// gen := combin.NewPermutationGenerator(len(g.chars), wordLen)
			// for gen.Next() {
			// 	select {
			// 	case <-ctx.Done():
			// 		return
			// 	default:
			// 	}
			// 	newPerm := gen.Permutation(nil)
			// 	newStr := permToStr(newPerm, g.chars)
			// 	if visited[string(newStr)] {
			// 		continue
			// 	}
			// 	visited[string(newStr)] = true
			// 	if !processPerm(ctx, space, wordLen, newStr, validChan, g.dictionary) {
			// 		return
			// 	}
			// }
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

// copies a board
func (b board) copy() (result board) {
	result = make(board, len(b))
	for i, row := range b {
		result[i] = make([]byte, len(row))
		copy(result[i], b[i])
	}
	return
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
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

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

// gets permutations of string a
// starts permuting at start and stops at the index before stop
// if all is true, the entire string is included, regardless of where stop is
// if all is false, the string up to stop
// func getPermutations(a string, start, stop int, all bool) (rt []string) {
// 	rt = make([]string, permLen(len(a), stop-start))
// 	i := 0
// 	strStop := stop
// 	strStart := start
// 	if all {
// 		strStop = len(a)
// 		strStart = 0
// 	}
// 	if stop < 0 {
// 		stop = len(a) + stop
// 	}
// 	perm([]byte(a), func(str []byte) {
// 		rt[i] = string(str[strStart:strStop])
// 		i++
// 	}, start, stop)
// 	fmt.Printf("num permutations: %d, should be: %d\n", len(rt), permLen(len(a), stop-start))
// 	return
// }

// Permute the values at index i to stop.
func perm(a []byte, f func([]byte), i, stop int, ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	if i > stop-1 {
		f(a)
		return
	}
	perm(a, f, i+1, stop, ctx)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1, stop, ctx)
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

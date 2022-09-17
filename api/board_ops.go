package function

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

// copies a board
func (b board) copy() (result board) {
	result = make(board, len(b))
	for i, row := range b {
		result[i] = make([]byte, len(row))
		copy(result[i], b[i])
	}
	return
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

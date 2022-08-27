package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DICTLEN int = 178691
const MAXLEN int = 7
const MINLEN int = 3

func main() {
	dictionary := getDictionary("dictionary.txt")
	fmt.Printf("dictionary length: %d\n", len(dictionary))
	userReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("enter characters: \n")
		chars, _ := userReader.ReadString('\n')
		chars = strings.ToUpper(strings.Replace(chars, "\n", "", -1))
		game := newGame(chars, "dictionary.txt")
		solution := game.solve()
		solution.print()
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestValidWord(t *testing.T) {
	dictionary := getDictionary("dictionary.txt")
	if !validWord("HELLO", dictionary) {
		t.Errorf("hello not found")
	}
	if validWord("ABCDEFG", dictionary) {
		t.Errorf("abcdef found in dictionary")
	}
	if !validWord("AA", dictionary) {
		t.Errorf("aa not found")
	}
	if !validWord("ZZZ", dictionary) {
		t.Errorf("ZZZ not found")
	}
}

func TestGetDictionary(t *testing.T) {
	dictionary := getDictionary("dictionary.txt")
	if len(dictionary) != DICTLEN {
		t.Errorf("failed to properly load dictionary")
	}
}

func TestPrefixMatches(t *testing.T) {
	dictionary := getDictionary("dictionary.txt")
	t.Logf("%v", findPrefixMatches("APP", dictionary))
}

func TestGetPermutations(t *testing.T) {
	//dictionary := getDictionary("dictionary.txt")
	perms := getPermutations("ab", 0, 2, false)
	t.Logf("%v", perms)
	perms = getPermutations("hel", 0, 2, true)
	t.Logf("%v", perms)
}

var testBoard1 board = [][]byte{
	{' ', ' ', 'A', ' '},
	{' ', ' ', 'B', ' '},
	{' ', ' ', 'C', ' '},
}

var testBoard2 board = [][]byte{
	{' ', 'R', 'A', ' '},
	{' ', ' ', 'B', ' '},
	{' ', ' ', 'C', ' '},
}

var testBoard3 board = [][]byte{
	{' ', ' ', 'Y', ' '},
	{' ', ' ', 'Z', ' '},
	{' ', 'R', 'A', ' '},
	{' ', ' ', 'B', 'D'},
	{' ', ' ', 'C', ' '},
}

// d gets picked up oops

func TestSpaceLeft(t *testing.T) {
	sl, _ := spaceLeft(0, 2, testBoard1)
	sl2, _ := spaceLeft(0, 2, testBoard2)
	if sl != -1 {
		t.Error()
	}
	if sl2 != 0 {
		t.Error()
	}
}

func TestSpaceRight(t *testing.T) {
	sl, _ := spaceRight(0, 2, testBoard1)
	sl2, _ := spaceRight(0, 1, testBoard2)
	if sl != -1 {
		t.Error()
	}
	if sl2 != 0 {
		t.Error()
	}
}

func TestGetSpaces(t *testing.T) {
	spaces1 := testBoard1.getHorizontalSpaces(false)
	spaces2 := testBoard2.getHorizontalSpaces(false)
	fmt.Printf("%v\n%v\n", spaces1, spaces2)
	if len(spaces1) != 3 || len(spaces2) != 2 {
		t.Error()
	}
	s1 := testBoard1.getSpaces()
	// s2 := testBoard2.getSpaces()
	s3 := testBoard3.getSpaces()
	fmt.Printf("%v\n%v\n", s1, s3)
}

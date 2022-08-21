package main

import "testing"

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

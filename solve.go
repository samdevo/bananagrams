package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func solveGame(chars string, dictionary []string) {
	firstWordOptions := findValidWords(chars, dictionary)
	fmt.Printf("%v\n", firstWordOptions)
}

func findValidWords(chars string, dictionary []string) (words []string) {
	if len(chars) < MINLEN {
		return
	}

	startLen := MAXLEN
	if len(chars) < MAXLEN {
		startLen = len(chars)
	}

	for wordLen := startLen; wordLen >= MINLEN; wordLen-- {
		fmt.Println(wordLen)
		combs := combinations([]byte(chars), wordLen)
		var perms [][]byte
		for _, comb := range combs {
			perms = append(perms, permutations(comb)...)
		}
		for _, word := range perms {
			if validWord(string(word), dictionary) {
				words = append(words, string(word))
			}
		}
	}
	return
}

func getDictionary(filename string) []string {
	var dictionary []string
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary = append(dictionary, scanner.Text())
	}
	return dictionary
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

func combinations(iterable []byte, r int) (rt [][]byte) {
	pool := iterable
	n := len(pool)

	if r > n {
		return
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]byte, r)
	for i, el := range indices {
		result[i] = pool[el]
	}
	s2 := make([]byte, r)
	copy(s2, result)
	rt = append(rt, s2)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i -= 1 {
		}

		if i < 0 {
			return
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		s2 = make([]byte, r)
		copy(s2, result)
		rt = append(rt, s2)
	}

}

// given a string, generate all arrangements of characters
func permutations(iterable []byte) (rt [][]byte) {
	n := len(iterable)
	if n == 1 {
		rt = append(rt, iterable)
		return
	}
	for i := 0; i < n; i++ {
		s1 := make([]byte, n-1)
		copy(s1, iterable[:i])
		copy(s1[i:], iterable[i+1:])
		perms := permutations(s1)
		for _, perm := range perms {
			rt = append(rt, append([]byte{iterable[i]}, perm...))
		}
	}
	return
}

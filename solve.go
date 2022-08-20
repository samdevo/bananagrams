package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func solveGame(chars string, dictionary []string) [][]byte {
	var board [][]byte
	// firstWordOptions := findValidWords(chars, dictionary)
	// fmt.Printf("%v\n", firstWordOptions)
	var wg sync.WaitGroup
	solution := make(chan [][]byte)
	fail := make(chan bool)
	go func() {
		wg.Add(1)
		search(&wg, chars, dictionary, board, solution)
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

func search(wg *sync.WaitGroup, chars string, dictionary []string, board [][]byte, done chan [][]byte) {
	defer wg.Done()
	ch := make(chan []byte)

}

func findValidWords(chars string, dictionary []string) (words []string) {
	if len(chars) < MINLEN {
		return
	}

	startLen := MAXLEN
	if len(chars) < MAXLEN {
		startLen = len(chars)
	}
	numValid := 0
	for wordLen := startLen; wordLen >= MINLEN; wordLen-- {
		perm([]byte(chars), func(str []byte) {
			if validWord(string(str[:wordLen]), dictionary) {
				words = append(words, string(str[:wordLen]))
				numValid++
			}
		}, 0, wordLen)
	}
	return
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

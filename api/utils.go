package function

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
)

// given a permutation of indices and the list of characters left, produce the given word
func permToStr(inds []int, chars []byte) []byte {
	newStr := make([]byte, len(inds))
	for i, ind := range inds {
		newStr[i] = chars[ind]
	}
	return newStr
}

func isChar(val byte) bool {
	return val != EMPTY
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

// Permute the values at index i to stop, calling f for each permutation.
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

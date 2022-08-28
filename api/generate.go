package function

import (
	"math/rand"
)

type CPop struct {
	char byte
	pop  int
}

var cpops = []CPop{
	{'A', 13},
	{'B', 3},
	{'C', 3},
	{'D', 6},
	{'E', 18},
	{'F', 3},
	{'G', 4},
	{'H', 3},
	{'I', 12},
	{'J', 2},
	{'K', 2},
	{'L', 5},
	{'M', 3},
	{'N', 8},
	{'O', 11},
	{'P', 3},
	{'Q', 2},
	{'R', 9},
	{'S', 6},
	{'T', 9},
	{'U', 6},
	{'V', 3},
	{'W', 3},
	{'X', 2},
	{'Y', 3},
	{'Z', 2},
}

func generateChars(numChars int) string {
	ccpops := make([]CPop, len(cpops))
	copy(ccpops, cpops)
	if numChars > 144 {
		return ""
	}
	chars := make([]byte, numChars)
	for i := 0; i < numChars; i++ {
		chars[i] = generateCharFromCPops(ccpops)
	}
	return string(chars)
}

// select n random chars from cpops, without replacement
func generateCharsFromCPops(n int, ccpops []CPop) string {
	chars := make([]byte, n)
	for i := 0; i < n; i++ {
		chars[i] = generateCharFromCPops(ccpops)
	}
	return string(chars)
}

// selects a random char from cpops, without replacement. returns the char
// subtracts the pop of the char from the total pop of the cpops
func generateCharFromCPops(cpops []CPop) byte {
	totalPop := 0
	for _, cpop := range cpops {
		totalPop += cpop.pop
	}
	randInt := rand.Intn(totalPop)
	for i := range cpops {
		randInt -= cpops[i].pop
		if randInt <= 0 {
			cpops[i].pop--
			return cpops[i].char
		}
	}
	return ' '
}

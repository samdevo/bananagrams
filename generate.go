package function

import "math/rand"

type CPop struct {
	char byte
	pop  int
}

var cpops = []CPop{
	{'A', 8167},
	{'B', 1492},
	{'C', 2782},
	{'D', 4253},
	{'E', 12702},
	{'F', 2228},
	{'G', 2015},
	{'H', 6094},
	{'I', 6966},
	{'J', 153},
	{'K', 772},
	{'L', 4025},
	{'M', 2406},
	{'N', 6749},
	{'O', 7507},
	{'P', 1929},
	{'Q', 95},
	{'R', 5987},
	{'S', 6327},
	{'T', 9056},
	{'U', 2758},
	{'V', 978},
	{'W', 2360},
	{'X', 150},
	{'Y', 1974},
	{'Z', 74},
}

func generateChars(numChars int) string {
	chars := make([]byte, numChars)
	for i := 0; i < numChars; i++ {
		chars[i] = generateCharFromCPops()
	}
	return string(chars)
}

func generateCharFromCPops() byte {
	total := 0
	for _, cpop := range cpops {
		total += cpop.pop
	}
	r := rand.Intn(total)
	for _, cpop := range cpops {
		r -= cpop.pop
		if r <= 0 {
			return cpop.char
		}
	}
	return 'A'
}

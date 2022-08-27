package function

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorld", helloWorld)
}

// helloWorld writes "Hello, World!" to the HTTP response.
func helloWorld(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}
	chars := strings.ToUpper(string(b))
	game := newGame(chars, "dictionary.txt")
	solution := game.solve()
	if solution != nil {
		solutionFlat := ""
		for _, line := range solution {
			solutionFlat += string(line) + "\n"
		}
		fmt.Fprintf(w, "%v\n", solutionFlat)
	} else {
		fmt.Fprintf(w, "No solution found\n")
	}
}

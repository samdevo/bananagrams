package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var mux = newMux()

func init() {
	functions.HTTP("HelloWorld", entryPoint)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/solve", solveHandler)
	mux.HandleFunc("/generate", generateHandler)
	// mux.HandleFunc("/subroute/three", three)

	return mux
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

type SolveRequest struct {
	Chars string `json:"chars"`
}

// helloWorld writes "Hello, World!" to the HTTP response.
func solveHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	bodyBytes := buf.Bytes()
	var req SolveRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chars := strings.ToUpper(string(req.Chars))
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

type GenerateRequest struct {
	NumChars int `json:"numChars"`
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	bodyBytes := buf.Bytes()
	var req GenerateRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chars := generateChars(req.NumChars)
	fmt.Fprintf(w, "%v\n", chars)
}

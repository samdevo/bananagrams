package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var mux = newMux()

const gcloudFuncSourceDir = "serverless_function_source_code"

func fixDir() {
	fileInfo, err := os.Stat(gcloudFuncSourceDir)
	if err == nil && fileInfo.IsDir() {
		_ = os.Chdir(gcloudFuncSourceDir)
	}
}

func init() {
	fixDir()
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
	boardStream := make(chan board)
	var wg sync.WaitGroup
	var solution board
	go func() {
		wg.Add(1)
		defer wg.Done()
		solution = game.solve(boardStream)
	}()
	for board := range boardStream {
		solutionFlat := ""
		for _, line := range board {
			solutionFlat += string(line) + "\n"
		}
		fmt.Fprintf(w, "%v", solutionFlat)
	}
	wg.Wait()
	if solution == nil {
		fmt.Fprintln(w, "Solution found!")
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

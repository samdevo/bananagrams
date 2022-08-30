package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mux.ServeHTTP(w, r)
}

type SolveRequest struct {
	Chars  string `json:"chars"`
	MinLen int    `json:"minChars"`
	MaxLen int    `json:"maxChars"`
}

// helloWorld writes "Hello, World!" to the HTTP response.
func solveHandler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	bodyBytes := buf.Bytes()
	var req SolveRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	chars := strings.ToUpper(string(req.Chars))
	game := newGame(chars, "dictionary.txt")
	// boardStream := make(chan board)
	var solution board
	// go func() {
	// 	wg.Add(1)
	// 	defer wg.Done()
	// 	solution = game.solve()
	// }()
	solution = game.solveSetLen(req.MinLen, req.MaxLen)
	solution.print()
	// for board := range boardStream {
	solutionFlat := ""
	for _, line := range solution {
		solutionFlat += string(line) + "\n"
	}
	fmt.Fprintf(w, "%v", solutionFlat)

	if solution == nil {
		fmt.Fprintln(w, "no solution found :(")
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
	if req.NumChars > 144 {
		http.Error(w, "Too many characters", http.StatusBadRequest)
		return
	}
	chars := generateChars(req.NumChars)
	fmt.Fprintf(w, "%v\n", chars)
}

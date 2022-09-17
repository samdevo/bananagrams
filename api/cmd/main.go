package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	_ "example.com/bananagrams"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "3000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

// func main() {
// 	dictionary := getDictionary("dictionary.txt")
// 	fmt.Printf("dictionary length: %d\n", len(dictionary))
// 	userReader := bufio.NewReader(os.Stdin)
// 	for {
// 		fmt.Printf("enter characters: \n")
// 		chars, _ := userReader.ReadString('\n')
// 		chars = strings.ToUpper(strings.Replace(chars, "\n", "", -1))
// 		game := newGame(chars, "dictionary.txt")
// 		solution := game.solve()
// 		solution.print()
// 	}
// }

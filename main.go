package main

import "fmt"

import "github.com/bmoxb/ikou/lexer"

func main() {
	tokens := lexer.Tokenise("hello ( 12.5 ); my comment\n(1254.1 world)")
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

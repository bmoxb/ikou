package main

import "fmt"

import "github.com/bmoxb/ikou/lexer"

func main() {
	tokens := lexer.Tokenise("hello ( 12.5 )")
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

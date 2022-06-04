package main

import "fmt"

import "github.com/bmoxb/ikou/lexer"

func main() {
	tokens, err := lexer.Tokenise("1.2abc hello if( 12.5 ) ; my comment\n(1254.1 world) ((15))")

	if err != nil {
		fmt.Println(err)
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
}

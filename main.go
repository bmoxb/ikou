package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bmoxb/ikou/lexer"
	"github.com/bmoxb/ikou/parser"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := stdin.ReadString('\n')

		tokens, err := lexer.Tokenise(line)

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, tok := range tokens {
			fmt.Printf("| %v\n", tok)
		}

		ast, err := parser.Parse(tokens)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("| %v\n", ast)
	}
}

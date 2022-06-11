package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bmoxb/ikou/lexer"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := stdin.ReadString('\n')

		tokens, err := lexer.Tokenise(line)

		if err != nil {
			fmt.Println(err)

		} else {
			for _, tok := range tokens {
				fmt.Printf("| %v\n", tok)
			}
		}

		fmt.Println()
	}
}

package lexer

import "fmt"

type Token struct {
	Type           TokenType
	Line           uint
	Position       uint
	OriginalString string
}

func (t Token) String() string {
	return fmt.Sprintf("token `%s` at line %d, position %d", t.OriginalString, t.Line, t.Position)
}

type TokenType uint

const (
	OpenTok       TokenType = iota // (
	CloseTok                       // )
	ColonTok                       // :
	CommentTok                     // ; my comment \n
	IntTok                         // 20
	FloatTok                       // 26.05
	IdentifierTok                  // identifier
	FunctionTok                    // fn
	LambdaTok                      // lambda
	IfTok                          // if
	LetTok                         // let
	DefineTok                      // define
)

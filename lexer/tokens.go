package lexer

import "fmt"

type Token struct {
	Type           TokenType
	Line           uint
	Position       uint
	OriginalString string
}

func (t Token) String() string {
	return fmt.Sprintf("%s token `%s` at line %d, position %d", t.Type, t.OriginalString, t.Line, t.Position)
}

type TokenType uint

const (
	OpenTok       TokenType = iota // (
	CloseTok                       // )
	ColonTok                       // :
	IntTok                         // 20
	FloatTok                       // 26.05
	IdentifierTok                  // identifier
	FunctionTok                    // fn
	LambdaTok                      // lambda
	IfTok                          // if
	LetTok                         // let
	DefineTok                      // define
)

func (t TokenType) String() string {
	switch t {
	case OpenTok:
		return "open bracket"
	case CloseTok:
		return "close bracket"
	case ColonTok:
		return "colon"
	case IntTok:
		return "integer"
	case FloatTok:
		return "float"
	case IdentifierTok:
		return "identifier"
	case FunctionTok:
		return "function keyword"
	case LambdaTok:
		return "lambda keyword"
	case IfTok:
		return "if keyword"
	case LetTok:
		return "let keyword"
	case DefineTok:
		return "define keyword"
	}
	return ""
}

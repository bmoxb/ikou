package lexer

import "fmt"

type Token struct {
	Type           TokenType
	Position       TokenPosition
	OriginalString string
}

func (t Token) String() string {
	return fmt.Sprintf("%s token %#v at %v", t.Type, t.OriginalString, t.Position)
}

type TokenType uint

const (
	OpenTok        TokenType = iota // (
	CloseTok                        // )
	SquareOpenTok                   // [
	SquareCloseTok                  // ]
	ColonTok                        // :
	QuoteTok                        // '
	BackquoteTok                    // ,
	IntTok                          // 20
	FloatTok                        // 26.05
	IdentifierTok                   // identifier
	StringTok                       // "Hello, world!"
	CharacterTok                    // \a
	TrueTok                         // true
	FalseTok                        // false
	LambdaTok                       // lambda
	IfTok                           // if
	LetTok                          // let
	DefineTok                       // define
)

func (t TokenType) String() string {
	switch t {
	case OpenTok:
		return "open ( bracket"
	case CloseTok:
		return "close ) bracket"
	case SquareOpenTok:
		return "square open [ bracket"
	case SquareCloseTok:
		return "square close ] bracket"
	case ColonTok:
		return "colon :"
	case QuoteTok:
		return "quote '"
	case BackquoteTok:
		return "backquote ,"
	case IntTok:
		return "integer literal"
	case FloatTok:
		return "float literal"
	case IdentifierTok:
		return "identifier"
	case StringTok:
		return "string literal"
	case CharacterTok:
		return "character literal"
	case TrueTok:
		return "true"
	case FalseTok:
		return "false"
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

type TokenPosition struct {
	Line               uint
	HorizontalPosition uint
}

func (t TokenPosition) String() string {
	return fmt.Sprintf("line %d, horizontal position %d", t.Line, t.HorizontalPosition)
}

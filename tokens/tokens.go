package tokens

import "fmt"

type Token struct {
	Type           Type
	Position       Position
	OriginalString string
}

func (t Token) String() string {
	return fmt.Sprintf("%s token %#v at %v", t.Type, t.OriginalString, t.Position)
}

type Type uint

const (
	Open        Type = iota // (
	Close                   // )
	SquareOpen              // [
	SquareClose             // ]
	Colon                   // :
	Quote                   // '
	Backquote               // ,
	Int                     // 20
	Float                   // 26.05
	Identifier              // identifier
	String                  // "Hello, world!"
	Character               // \a
	True                    // true
	False                   // false
	Lambda                  // lambda
	If                      // if
	Let                     // let
	Define                  // define
)

func (t Type) String() string {
	switch t {
	case Open:
		return "open ( bracket"
	case Close:
		return "close ) bracket"
	case SquareOpen:
		return "square open [ bracket"
	case SquareClose:
		return "square close ] bracket"
	case Colon:
		return "colon :"
	case Quote:
		return "quote '"
	case Backquote:
		return "backquote ,"
	case Int:
		return "integer literal"
	case Float:
		return "float literal"
	case Identifier:
		return "identifier"
	case String:
		return "string literal"
	case Character:
		return "character literal"
	case True:
		return "true"
	case False:
		return "false"
	case Lambda:
		return "lambda keyword"
	case If:
		return "if keyword"
	case Let:
		return "let keyword"
	case Define:
		return "define keyword"
	}
	return ""
}

type Position struct {
	Line               uint
	HorizontalPosition uint
}

func (t Position) String() string {
	return fmt.Sprintf("line %d, horizontal position %d", t.Line, t.HorizontalPosition)
}

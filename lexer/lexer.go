package lexer

import (
	//"log"
	"fmt"
	"strings"
)

var keywordMap = map[string]TokenType{
	"true":   TrueTok,
	"false":  FalseTok,
	"lambda": LambdaTok,
	"if":     IfTok,
	"let":    LetTok,
	"define": DefineTok,
}

var finalTokenTypeMap = map[state]TokenType{
	identState: IdentifierTok,
	intState:   IntTok,
	floatState: FloatTok,
}

var singleCharacterTokenMap = map[rune]TokenType{
	'(':  OpenTok,
	')':  CloseTok,
	'[':  SquareOpenTok,
	']':  SquareCloseTok,
	':':  ColonTok,
	'\'': QuoteTok,
	',':  BackquoteTok,
}

func Tokenise(input string) ([]Token, error) {
	if input == "" {
		return make([]Token, 0), nil
	}

	l := lexer{
		currentState: initialState,
		pos:          TokenPosition{Line: 1, HorizontalPosition: 0},
	}

	lines := strings.Split(input, "\n")
	linesIndex := 0

	runes := []rune(input)
	c := runes[0]

	for _, peek := range runes[1:] {
		if c == '\n' {
			linesIndex += 1
		}

		err := l.processChar(c, peek, lines[linesIndex])
		if err != nil {
			return nil, err
		}

		c = peek
	}

	err := l.processChar(c, 0, lines[linesIndex])
	if err != nil {
		return nil, err
	}

	l.eof()

	return l.tokens, nil
}

type lexer struct {
	tokens        []Token
	currentState  state
	currentString strings.Builder
	pos           TokenPosition
}

func (l *lexer) processChar(c, peek rune, currentLine string) error {
	//log.Printf("processing character '%c' (peek '%c') in state %d with current string \"%s\"", c, peek, l.currentState, l.currentString.String())

	if c == '\n' {
		l.pos.Line += 1
		l.pos.HorizontalPosition = 0

		if l.currentState == commentState {
			l.discardToken()
		}

		return nil
	}

	l.pos.HorizontalPosition += 1
	l.currentString.WriteRune(c)

	switch l.currentState {
	case initialState:
		singleCharTok, isSingleCharTok := singleCharacterTokenMap[c]

		if isSingleCharTok {
			l.addToken(singleCharTok)

		} else if c == ';' {
			// A semicolon ; character indicates the start of a comment so change to
			// comment state.

			l.currentState = commentState

		} else if c == '"' {
			// A double quote " character indicates the start of a string so change to
			// string state.

			l.currentState = stringState

		} else if runeIsNumeral(c) || (c == '~' && runeIsNumeral(peek)) {
			// If character is `[0-9]` or character and peek are `~[0-9]` (i.e., negative
			// number)...

			// If more than just the one numeral then change state, otherwise add the
			// single number as an integer token.
			if runeIsNumeral(peek) || peek == '.' {
				l.currentState = intState
			} else {
				l.addToken(IntTok)
			}

		} else if c == '.' && runeIsNumeral(peek) {
			// If character and peek are in the form `\.[0-9]` (i.e., a float literal),
			// change to float state.

			l.currentState = floatState

		} else if runeIsIdentChar(c) {
			// If character could be part of an identifier...

			// If peek is also an identifier character then change state.
			if runeIsIdentChar(peek) {
				l.currentState = identState
				// Ensure the identifier is not immediately followed by a number literal.
			} else if runeIsNumeral(peek) || runeIsOneOf(peek, ".-") {
				return &LexicalError{pos: l.pos, line: currentLine, msg: "identifier and number literal must be separated by whitespace"}
				// If there is just a single identifier character alone then add it is as an identifier token.
			} else {
				l.addToken(IdentifierTok)
			}

		} else {
			if runeIsOneOf(c, " \t\n\r") {
				// Discard any whitespace characters.

				l.discardToken()

			} else {
				// Any non-whitespace unexpected characters result in an error.

				l.pos.HorizontalPosition -= 1
				return &LexicalError{pos: l.pos, line: currentLine, msg: fmt.Sprintf("unexpected character: %q", c)}
			}
		}

	case identState:
		if peek == '.' {
			return &LexicalError{pos: l.pos, line: currentLine, msg: "identifier and floating-point literal must be separated by whitespace"}
		}

		if !runeIsIdentChar(peek) {
			keyword, isKeyword := keywordMap[l.currentString.String()]

			if isKeyword {
				l.addToken(keyword)
			} else {
				l.addToken(IdentifierTok)
			}
		}

	case intState:
		if !runeIsNumeral(peek) && runeIsIdentChar(peek) {
			return &LexicalError{pos: l.pos, line: currentLine, msg: "integer literal and identifier must be separated by whitespace"}
		}

		if c == '.' {
			l.currentState = floatState

		} else if peek != '.' && !runeIsNumeral(peek) {
			l.addToken(IntTok)
		}

	case floatState:
		if peek == '.' {
			return &LexicalError{pos: l.pos, line: currentLine, msg: "found multiple decimal point characters found in floating-point literal"}

		}

		if !runeIsNumeral(peek) {
			if runeIsIdentChar(peek) {
				return &LexicalError{pos: l.pos, line: currentLine, msg: "floating-point literal and identifier must be separated by whitespace"}
			}

			l.addToken(FloatTok)
		}

	case stringState:
		if c == '"' {
			l.addToken(StringTok)
		}

		if c == '\\' && !runeIsOneOf(peek, `"tn\`) {
			return &LexicalError{pos: l.pos, line: currentLine, msg: fmt.Sprintf("found invalid escape sequence character %#v found in string literal", peek)}
		}
	}

	return nil
}

func (l *lexer) eof() {
	ty, isFinalToken := finalTokenTypeMap[l.currentState]
	if isFinalToken {
		l.addToken(ty)
	}
}

func (l *lexer) addToken(ty TokenType) {
	tok := Token{
		Type:           ty,
		Position:       l.pos,
		OriginalString: l.currentString.String(),
	}

	l.currentState = initialState
	l.currentString.Reset()

	l.tokens = append(l.tokens, tok)
}

func (l *lexer) discardToken() {
	l.currentState = initialState
	l.currentString.Reset()
}

type state int

const (
	initialState state = iota
	identState
	commentState
	intState
	floatState
	stringState
)

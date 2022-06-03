package lexer

import (
	//"log"
	"fmt"
	"strings"
)

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
		switch c {
		case '(':
			l.addToken(OpenTok)
		case ')':
			l.addToken(CloseTok)
		case ':':
			l.addToken(ColonTok)
		case ';':
			l.currentState = commentState
		default:
			if runeIsNumeral(c) || (c == '-' && runeIsNumeral(peek)) {
				if runeIsNumeral(peek) || peek == '.' {
					l.currentState = intState
				} else {
					l.addToken(IntTok)
				}
			} else if c == '.' && runeIsNumeral(peek) {
				l.currentState = floatState
			} else if runeIsIdentChar(c) {
				if runeIsIdentChar(peek) {
					l.currentState = identState
				} else {
					l.addToken(IdentifierTok)
				}
			} else {
				if runeIsOneOf(c, " \t\n\r") {
					l.discardToken()
				} else {
					l.pos.HorizontalPosition -= 1
					return &LexicalError{pos: l.pos, line: currentLine, msg: fmt.Sprintf("unexpected character: %q", c)}
				}
			}
		}

	case identState:
		if !runeIsIdentChar(peek) {
			keywords := map[string]TokenType{
				"fn":     FunctionTok,
				"lambda": LambdaTok,
				"if":     IfTok,
				"let":    LetTok,
				"define": DefineTok,
			}

			keyword, isKeyword := keywords[l.currentString.String()]

			if isKeyword {
				l.addToken(keyword)
			} else {
				l.addToken(IdentifierTok)
			}
		}

	case intState:
		if c == '.' {
			l.currentState = floatState
		} else if peek != '.' && !runeIsNumeral(peek) {
			l.addToken(IntTok)
		}

	case floatState:
		if peek == '.' {
			return &LexicalError{pos: l.pos, line: currentLine, msg: fmt.Sprintf("found multiple decimal point characters found in floating-point literal")}
		} else if !runeIsNumeral(peek) {
			l.addToken(FloatTok)
		}
	}

	return nil
}

func (l *lexer) eof() {
	finalTokenTypeMap := map[state]TokenType{
		identState: IdentifierTok,
		intState:   IntTok,
		floatState: FloatTok,
	}

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
)

package lexer

import (
	//"log"
	"errors"
	"fmt"
	"strings"
)

func Tokenise(input string) ([]Token, error) {
	if input == "" {
		return make([]Token, 0), nil
	}

	l := lexer{
		currentState: initialState,
		line:         1,
		pos:          0,
	}

	runes := []rune(input)
	c := runes[0]

	for _, peek := range runes[1:] {
		err := l.processChar(c, peek)
		if err != nil {
			return nil, err
		}

		c = peek
	}

	err := l.processChar(c, 0)
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
	line          uint
	pos           uint
}

func (l *lexer) processChar(c, peek rune) error {
	//log.Printf("processing character '%c' (peek '%c') in state %d with current string \"%s\"", c, peek, l.currentState, l.currentString.String())

	if c == '\n' {
		l.line += 1
		l.pos = 0

		if l.currentState == commentState {
			l.discardToken()
		}

		return nil
	}

	l.pos += 1
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
					return errors.New(fmt.Sprintf("unexpected character: %q", c))
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
		if !runeIsNumeral(peek) {
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
		Line:           l.line,
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

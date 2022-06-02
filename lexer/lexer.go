package lexer

import (
	//"log"
	"strings"
)

func Tokenise(input string) []Token {
	if input == "" {
		return make([]Token, 0)
	}

	l := lexer{
		currentState: initialState,
		line:         1,
		pos:          0,
	}

    runes := []rune(input)
	c := runes[0]
	for _, peek := range runes[1:] {
		l.processChar(c, peek)
		c = peek
	}
	l.processChar(c, 0)
    l.eof()

	return l.tokens
}

type lexer struct {
	tokens        []Token
	currentState  state
	currentString strings.Builder
	line          uint
	pos           uint
}

func (l *lexer) processChar(c, peek rune) {
	//log.Printf("processing character '%c' (peek '%c') in state %d with current string \"%s\"", c, peek, l.currentState, l.currentString.String())

	if c == '\n' {
		l.line += 1
		l.pos = 0

		if l.currentState == commentState {
			l.discardToken()
		}

		return
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
			if runeIsNumeral(c) {
				l.currentState = intState
			} else if runeIsIdentChar(c) {
				l.currentState = identState
			} else {
				l.discardToken()
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
}

func (l *lexer) eof() {
    finalTokenTypeMap := map[state]TokenType{
        identState: IdentifierTok,
        intState: IntTok,
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

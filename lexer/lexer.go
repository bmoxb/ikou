package lexer

import (
	//"log"
	"fmt"
	"strings"
)

var keywordTokenTypeMap = map[string]TokenType{
	"true":   TrueTok,
	"false":  FalseTok,
	"lambda": LambdaTok,
	"if":     IfTok,
	"let":    LetTok,
	"define": DefineTok,
}

var singleCharacterTokenTypeMap = map[rune]TokenType{
	'(':  OpenTok,
	')':  CloseTok,
	'[':  SquareOpenTok,
	']':  SquareCloseTok,
	':':  ColonTok,
	'\'': QuoteTok,
	',':  BackquoteTok,
}

var validCharacterLiterals = map[string]struct{}{
	`\space`:   {},
	`\tab`:     {},
	`\newline`: {},
}

var whitespaceNameMap = map[rune]string{
	' ':  "space",
	'\t': "tab",
	'\n': "newline",
	'\r': "newline",
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

	err = l.eof(lines[linesIndex])
	if err != nil {
		return nil, err
	}

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
		ty, isSingleCharTok := singleCharacterTokenTypeMap[c]

		if isSingleCharTok {
			l.addSpecificTokenType(ty)

		} else if c == ';' {
			// A semicolon ; character indicates the start of a comment so change to
			// comment state.

			l.currentState = commentState

		} else if c == '"' {
			// A double quote " character indicates the start of a string so change to
			// string state.

			l.currentState = stringState

		} else if c == '\\' {
			// A backslash \ character indicates a character literal so change to
			// character state.

			if runeIsWhitespace(peek) {
				whitespace := whitespaceNameMap[peek]
				return l.makeError(currentLine, fmt.Sprintf("A %s character cannot be used as a character literal - consider using `\\%v` instead", whitespace, whitespace))
			}

			l.currentState = characterState

		} else if runeIsNumeral(c) || (c == '~' && runeIsNumeral(peek)) {
			// If character is `[0-9]` or character and peek are `~[0-9]` (i.e., negative
			// number)...

			if !runeIsNumeral(peek) && runeIsIdentChar(peek) {
				return l.makeIntIdentSepError(currentLine)
			}

			l.currentState = intState

			// If the following token is not a numeral nor a decimal point then just add
			// the token.
			if !runeIsNumeral(peek) && peek != '.' {
				return l.addToken(currentLine)
			}

		} else if c == '.' && runeIsNumeral(peek) {
			// If character and peek are in the form `\.[0-9]` (i.e., a float literal),
			// change to float state.

			l.currentState = floatState

		} else if runeIsIdentChar(c) {
			// If character could be part of an identifier...

			if peek == '.' {
				return l.makeIdentPeriodSepError(currentLine)
			}

			l.currentState = identState

			if !runeIsIdentChar(peek) {
				// If there is just a single identifier character alone then add it is as an identifier token.
				return l.addToken(currentLine)
			}

		} else {
			if runeIsWhitespace(c) {
				// Discard any whitespace characters.

				l.discardToken()

			} else {
				// Any non-whitespace unexpected characters result in an error.

				l.pos.HorizontalPosition -= 1
				return l.makeError(currentLine, fmt.Sprintf("unexpected character: %q", c))
			}
		}

	case identState:
		if c == '.' || peek == '.' {
			return l.makeIdentPeriodSepError(currentLine)
		}

		if !runeIsIdentChar(peek) {
			return l.addToken(currentLine)
		}

	case intState:
		if !runeIsNumeral(peek) && runeIsIdentChar(peek) {
			return l.makeIntIdentSepError(currentLine)
		}

		if c == '.' {
			l.currentState = floatState

		} else if peek != '.' && !runeIsNumeral(peek) {
			return l.addToken(currentLine)
		}

	case floatState:
		if peek == '.' {
			return l.makeError(currentLine, "found multiple decimal point characters found in floating-point literal")
		}

		if !runeIsNumeral(peek) {
			if runeIsIdentChar(peek) {
				return l.makeError(currentLine, "floating-point literal and identifier must be separated by whitespace")
			}

			return l.addToken(currentLine)
		}

	case stringState:
		if c == '"' {
			return l.addToken(currentLine)

		} else if c == '\\' && !runeIsOneOf(peek, `"tn\`) {
			return l.makeError(currentLine, fmt.Sprintf("found invalid escape sequence \\%c found in string literal", peek))
		}

	case characterState:
		if !runeIsAlpha(peek) {
			return l.addToken(currentLine)
		}
	}

	return nil
}

func (l *lexer) eof(currentLine string) error {
	if l.currentState == stringState {
		return l.makeError(currentLine, "reached end of file during evaluation of a string literal")
	}

	return l.addToken(currentLine)
}

func (l *lexer) addToken(currentLine string) error {
	var ty TokenType

	switch l.currentState {
	case initialState, commentState:
		return nil // discard

	case identState:
		keywordTokType, isKeyword := keywordTokenTypeMap[l.currentString.String()]

		if isKeyword {
			ty = keywordTokType
		} else {
			ty = IdentifierTok
		}

	case intState:
		ty = IntTok

	case floatState:
		ty = FloatTok

	case stringState:
		ty = StringTok

	case characterState:
		ty = CharacterTok

		s := l.currentString.String()

		_, valid := validCharacterLiterals[s]

		if len(s) < 2 || (len(s) > 2 && !valid) {
			return l.makeError(currentLine, "invalid character literal")
		}
	}

	l.addSpecificTokenType(ty)

	return nil
}

func (l *lexer) addSpecificTokenType(ty TokenType) {
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

func (l *lexer) makeError(currentLine, msg string) error {
	return &LexicalError{pos: l.pos, line: currentLine, msg: msg}
}

func (l *lexer) makeIdentPeriodSepError(currentLine string) error {
	return l.makeError(currentLine, "identifier and '.' must be separated by whitespace")
}

func (l *lexer) makeIntIdentSepError(currentLine string) error {
	return l.makeError(currentLine, "integer literal and identifier must be separated by whitespace")
}

type state int

const (
	initialState state = iota
	identState
	commentState
	intState
	floatState
	stringState
	characterState
)

package lexer

import "testing"

type tokenTest struct {
	ty TokenType
	s  string
}

func tt(ty TokenType, s string) tokenTest {
	return tokenTest{ty, s}
}

func tp(line, pos uint) TokenPosition {
	return TokenPosition{Line: line, HorizontalPosition: pos}
}

func TestTokeniseSingleTokens(t *testing.T) {
	table := map[string]tokenTest{
		// Single character tokens:
		" (":    tt(OpenTok, "("),
		") ":    tt(CloseTok, ")"),
		"[":     tt(SquareOpenTok, "["),
		" ] ":   tt(SquareCloseTok, "]"),
		"\t: ":  tt(ColonTok, ":"),
		"'":     tt(QuoteTok, "'"),
		"\n,\n": tt(BackquoteTok, ","),

		// Identifiers:
		"abc\n":        tt(IdentifierTok, "abc"),
		"abc-def+":     tt(IdentifierTok, "abc-def+"),
		"ABC_123":      tt(IdentifierTok, "ABC_123"),
		"\n+\n":        tt(IdentifierTok, "+"),
		"  -  ":        tt(IdentifierTok, "-"),
		"  +  ":        tt(IdentifierTok, "+"),
		"\t--  ":       tt(IdentifierTok, "--"),
		" a":           tt(IdentifierTok, "a"),
		"; comment\na": tt(IdentifierTok, "a"),
		"a-":           tt(IdentifierTok, "a-"),
		"aa-":          tt(IdentifierTok, "aa-"),
		"a-5":          tt(IdentifierTok, "a-5"),
		"-1":           tt(IdentifierTok, "-1"),
		"-":            tt(IdentifierTok, "-"),

		// Numbers:
		"12":    tt(IntTok, "12"),
		"0":     tt(IntTok, "0"),
		"1.":    tt(FloatTok, "1."),
		".1":    tt(FloatTok, ".1"),
		"~1234": tt(IntTok, "~1234"),
		"12.5":  tt(FloatTok, "12.5"),
		"~0.5":  tt(FloatTok, "~0.5"),

		// Keywords:
		"let": tt(LetTok, "let"),
		"Let": tt(IdentifierTok, "Let"),
		"if":  tt(IfTok, "if"),
		"iff": tt(IdentifierTok, "iff"),

		// Strings:
		` "" `:             tt(StringTok, `""`),
		`"abc def"`:        tt(StringTok, `"abc def"`),
		`"\tHello, 世界！\n"`: tt(StringTok, `"\tHello, 世界！\n"`),
	}

	for input, expected := range table {
		tokens, err := Tokenise(input)

		if err != nil {
			t.Errorf("Tokenise(%#v) produced unexpected error: \n%v", input, err)
		} else if len(tokens) != 1 {
			t.Errorf("Tokenise(%#v) expected to return 1 token but instead returned %d tokens", input, len(tokens))
		} else {
			tok := tokens[0]

			if tok.Type != expected.ty {
				t.Errorf("Tokenise(%#v) gave token type %v but expected %v", input, tok.Type, expected.ty)
			}

			if tok.OriginalString != expected.s {
				t.Errorf("Tokenise(%#v) gave token string %#v but expected %#v", input, tok.OriginalString, expected.s)
			}
		}
	}

	invalid := []string{
		"1.2.1",
		".",
		"@",
		"a\\",
		"12abc",
		"5.5abc",
		"5.a",
		"55.aa",
		"a.5",
		"aa.55",
		"a.a",
		"aa.aa",
		".a",
		"a.",
		"aa-12.5",
		"a#",
		"#",
		"aa#",
		"#aa",
		"~",
		"a~",
	}

	for _, input := range invalid {
		_, err := Tokenise(input)

		if err == nil {
			t.Errorf("Input %#v is invalid but no error was returned during lexing", input)
		}
	}
}

func TestTokeniseMultipleTokens(t *testing.T) {
	table := map[string][]tokenTest{
		// Input edge cases:
		"":   []tokenTest{},
		" ":  []tokenTest{},
		";":  []tokenTest{},
		"\n": []tokenTest{},

		"(+ 15 25)":                           []tokenTest{tt(OpenTok, "("), tt(IdentifierTok, "+"), tt(IntTok, "15"), tt(IntTok, "25"), tt(CloseTok, ")")},
		"\t( - ~0.1 0.2 )\n":                  []tokenTest{tt(OpenTok, "("), tt(IdentifierTok, "-"), tt(FloatTok, "~0.1"), tt(FloatTok, "0.2"), tt(CloseTok, ")")},
		"; comment\nlet0 0 let\tLET; comment": []tokenTest{tt(IdentifierTok, "let0"), tt(IntTok, "0"), tt(LetTok, "let"), tt(IdentifierTok, "LET")},
	}

	for input, expected := range table {
		tokens, err := Tokenise(input)

		if err != nil {
			t.Errorf("Tokenise(%#v) produced unexpected error: \n%v", input, err)
		} else if len(tokens) != len(expected) {
			t.Errorf("Tokenise(%#v) returned %d tokens but %d were expected", input, len(tokens), len(expected))
		} else {
			for i, tok := range tokens {
				if (tok.Type != expected[i].ty) || (tok.OriginalString != expected[i].s) {
					t.Errorf("Tokenise(%#v) returned %v where a token of type %v with string %#v was expected", input, tok, expected[i].ty, expected[i].s)
				}
			}
		}
	}
}

func TestTokenisePositionLine(t *testing.T) {
	table := map[string][]TokenPosition{
		"( 21 )":           []TokenPosition{tp(1, 1), tp(1, 4), tp(1, 6)},
		"12.5\n10":         []TokenPosition{tp(1, 4), tp(2, 2)},
		"; comment\nabc\n": []TokenPosition{tp(2, 3)},
	}

	for input, positions := range table {
		tokens, err := Tokenise(input)

		if err != nil {
			t.Errorf("Tokenise(%#v) produced unexpected error: \n%v", input, err)
		} else if len(tokens) != len(positions) {
			t.Errorf("Tokenise(%#v) returned %d tokens but %d were expected", input, len(tokens), len(positions))
		} else {
			for i, tok := range tokens {
				if tok.Position != positions[i] {
					t.Errorf("Tokenise(%#v) gave %v which was expected to be at %v", input, tok, positions[i])
				}
			}
		}
	}
}

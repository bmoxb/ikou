package lexer

import "testing"

type tokenTest struct {
	ty TokenType
	s  string
}

func tt(ty TokenType, s string) tokenTest {
	return tokenTest{ty, s}
}

func TestTokeniseSingleTokens(t *testing.T) {
	table := map[string]tokenTest{
		" (":           tt(OpenTok, "("),
		") ":           tt(CloseTok, ")"),
		"\t: ":         tt(ColonTok, ":"),
		"abc\n":        tt(IdentifierTok, "abc"),
		"abc-def+":     tt(IdentifierTok, "abc-def+"),
		"ABC_123":      tt(IdentifierTok, "ABC_123"),
		"\n+\n":        tt(IdentifierTok, "+"),
		"  -  ":        tt(IdentifierTok, "-"),
		"  +  ":        tt(IdentifierTok, "+"),
		"\t--  ":       tt(IdentifierTok, "--"),
		" a\\":         tt(IdentifierTok, "a"),
		"; comment\na": tt(IdentifierTok, "a"),
		"12":           tt(IntTok, "12"),
		"0":            tt(IntTok, "0"),
		"-1234":        tt(IntTok, "-1234"),
		"12.5":         tt(FloatTok, "12.5"),
		"-0.5":         tt(FloatTok, "-0.5"),
		"let":          tt(LetTok, "let"),
		"Let":          tt(IdentifierTok, "Let"),
		"if":           tt(IfTok, "if"),
		"iff":          tt(IdentifierTok, "iff"),
	}

	for input, expected := range table {
		tokens := Tokenise(input)

		if len(tokens) != 1 {
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
}

func TestTokeniseMultipleTokens(t *testing.T) {
	table := map[string][]tokenTest{
		"(+ 15 25)":                           []tokenTest{tt(OpenTok, "("), tt(IdentifierTok, "+"), tt(IntTok, "15"), tt(IntTok, "25"), tt(CloseTok, ")")},
		"\t( - -0.1 0.2 )\n":                  []tokenTest{tt(OpenTok, "("), tt(IdentifierTok, "-"), tt(FloatTok, "-0.1"), tt(FloatTok, "0.2"), tt(CloseTok, ")")},
		"; comment\nlet0 0 let\tLET; comment": []tokenTest{tt(IdentifierTok, "let0"), tt(IntTok, "0"), tt(LetTok, "let"), tt(IdentifierTok, "LET")},
	}

	for input, expected := range table {
		tokens := Tokenise(input)

		if len(tokens) != len(expected) {
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
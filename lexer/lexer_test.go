package lexer

import "testing"

type tokenTest struct {
	ty TokenType
	s  string
}

func TestTokeniseSingleTokens(t *testing.T) {
	table := map[string]tokenTest{
		" (":           tokenTest{ty: OpenTok, s: "("},
		") ":           tokenTest{ty: CloseTok, s: ")"},
		"\t: ":         tokenTest{ty: ColonTok, s: ":"},
		"abc\n":        tokenTest{ty: IdentifierTok, s: "abc"},
		"abc-def+":     tokenTest{ty: IdentifierTok, s: "abc-def+"},
		"ABC_123":      tokenTest{ty: IdentifierTok, s: "ABC_123"},
		"\n+\n":        tokenTest{ty: IdentifierTok, s: "+"},
		"  -  ":        tokenTest{ty: IdentifierTok, s: "-"},
		"  +  ":        tokenTest{ty: IdentifierTok, s: "+"},
		"\t--  ":       tokenTest{ty: IdentifierTok, s: "--"},
		" a\\":         tokenTest{ty: IdentifierTok, s: "a"},
		"; comment\na": tokenTest{ty: IdentifierTok, s: "a"},
		"12":           tokenTest{ty: IntTok, s: "12"},
		"0":            tokenTest{ty: IntTok, s: "0"},
		"-1234":        tokenTest{ty: IntTok, s: "-1234"},
		"12.5":         tokenTest{ty: FloatTok, s: "12.5"},
		"-0.5":         tokenTest{ty: FloatTok, s: "-0.5"},
		"let":          tokenTest{ty: LetTok, s: "let"},
		"Let":          tokenTest{ty: IdentifierTok, s: "Let"},
		"if":           tokenTest{ty: IfTok, s: "if"},
		"iff":          tokenTest{ty: IdentifierTok, s: "iff"},
	}

	for input, expected := range table {
		tokens := Tokenise(input)

		if len(tokens) != 1 {
			t.Errorf("Tokenise(\"%s\") expected to return 1 token but instead returned %d tokens", input, len(tokens))
		} else {
			tok := tokens[0]

			if tok.Type != expected.ty {
				t.Errorf("Tokenise(\"%#v\") gave token type %v but expected %v", input, tok.Type, expected.ty)
			}

			if tok.OriginalString != expected.s {
				t.Errorf("Tokenise(%#v) gave token string %#v but expected %#v", input, tok.OriginalString, expected.s)
			}
		}
	}
}

func TestTokeniseMultipleTokens(t *testing.T) {
}

package lexer

import (
	"testing"

	"github.com/bmoxb/ikou/tokens"
)

type tokenTest struct {
	ty tokens.Type
	s  string
}

func tt(ty tokens.Type, s string) tokenTest {
	return tokenTest{ty, s}
}

func tp(line, pos uint) tokens.Position {
	return tokens.Position{Line: line, HorizontalPosition: pos}
}

func TestTokeniseSingleTokens(t *testing.T) {
	table := map[string]tokenTest{
		// Single character tokens:
		" (":    tt(tokens.Open, "("),
		") ":    tt(tokens.Close, ")"),
		"[":     tt(tokens.SquareOpen, "["),
		" ] ":   tt(tokens.SquareClose, "]"),
		"\t: ":  tt(tokens.Colon, ":"),
		"'":     tt(tokens.Quote, "'"),
		"\n,\n": tt(tokens.Backquote, ","),

		// Identifiers:
		"abc\n":        tt(tokens.Identifier, "abc"),
		"abc-def+":     tt(tokens.Identifier, "abc-def+"),
		"ABC_123":      tt(tokens.Identifier, "ABC_123"),
		"\n+\n":        tt(tokens.Identifier, "+"),
		"  -  ":        tt(tokens.Identifier, "-"),
		"  +  ":        tt(tokens.Identifier, "+"),
		"\t--  ":       tt(tokens.Identifier, "--"),
		" a":           tt(tokens.Identifier, "a"),
		"; comment\na": tt(tokens.Identifier, "a"),
		"a-":           tt(tokens.Identifier, "a-"),
		"aa-":          tt(tokens.Identifier, "aa-"),
		"a-5":          tt(tokens.Identifier, "a-5"),
		"-1":           tt(tokens.Identifier, "-1"),
		"-":            tt(tokens.Identifier, "-"),
		"a5":           tt(tokens.Identifier, "a5"),

		// Numbers:
		"12":    tt(tokens.Int, "12"),
		"0":     tt(tokens.Int, "0"),
		"1.":    tt(tokens.Float, "1."),
		".1":    tt(tokens.Float, ".1"),
		"~1234": tt(tokens.Int, "~1234"),
		"12.5":  tt(tokens.Float, "12.5"),
		"~0.5":  tt(tokens.Float, "~0.5"),

		// Keywords:
		"let": tt(tokens.Let, "let"),
		"Let": tt(tokens.Identifier, "Let"),
		"if":  tt(tokens.If, "if"),
		"iff": tt(tokens.Identifier, "iff"),

		// Strings:
		` "" `:             tt(tokens.String, `""`),
		`"abc def"`:        tt(tokens.String, `"abc def"`),
		`"\tHello, 世界！\n"`: tt(tokens.String, `"\tHello, 世界！\n"`),

		// Characters:
		`\a`:       tt(tokens.Character, `\a`),
		`\tab`:     tt(tokens.Character, `\tab`),
		`\space`:   tt(tokens.Character, `\space`),
		`\newline`: tt(tokens.Character, `\newline`),
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
		// Multiple decimal points in float literal:
		"1.2.1",
		".2.",
		"1.2.",
		".2.1",

		// Identifiers and number literals:
		"5.5abc",
		"5.a",
		"55.aa",
		"a.5",
		"aa.55",
		"a0.5",
		"a.a",
		"aa.aa",
		".a",
		"a.",
		"12abc",
		"1a",

		// Negation symbol:
		"~",
		"a~",
		"~a",

		// Character literal:
		`a\`,
		`\`,
		` \ `,
		`\abc`,
		"\\\n",

		// EOF during string literal:
		`"`,
		`"abc`,

		// Misc:
		".",
		"@",
		"aa-12.5",
		"a#",
		"#",
		"aa#",
		"#aa",
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
		"":   {},
		" ":  {},
		";":  {},
		"\n": {},

		// Spacing:
		"\t( - ~0.1 0.2 )\n": {tt(tokens.Open, "("), tt(tokens.Identifier, "-"), tt(tokens.Float, "~0.1"), tt(tokens.Float, "0.2"), tt(tokens.Close, ")")},
		`( ( a "b" ) 12.5 )`: {tt(tokens.Open, "("), tt(tokens.Open, "("), tt(tokens.Identifier, "a"), tt(tokens.String, `"b"`), tt(tokens.Close, ")"), tt(tokens.Float, "12.5"), tt(tokens.Close, ")")},

		// Minimal spacing:
		"(+ 15 25)":                          {tt(tokens.Open, "("), tt(tokens.Identifier, "+"), tt(tokens.Int, "15"), tt(tokens.Int, "25"), tt(tokens.Close, ")")},
		`"abc""def"`:                         {tt(tokens.String, `"abc"`), tt(tokens.String, `"def"`)},
		"; comment\nlet0 0 let\tLET;comment": {tt(tokens.Identifier, "let0"), tt(tokens.Int, "0"), tt(tokens.Let, "let"), tt(tokens.Identifier, "LET")},
		"((\n12.5\n\"\"))":                   {tt(tokens.Open, "("), tt(tokens.Open, "("), tt(tokens.Float, "12.5"), tt(tokens.String, `""`), tt(tokens.Close, ")"), tt(tokens.Close, ")")},
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
	table := map[string][]tokens.Position{
		"( 21 )":           {tp(1, 1), tp(1, 4), tp(1, 6)},
		"12.5\n10":         {tp(1, 4), tp(2, 2)},
		"; comment\nabc\n": {tp(2, 3)},
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

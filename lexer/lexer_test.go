package lexer

import "testing"

func TestTokeniseSingleTokens(t *testing.T) {
    table := map[string]TokenType{
        "(": OpenTok,
        ")": CloseTok,
        ":": ColonTok,
        "abc": IdentifierTok,
        "abc-def+": IdentifierTok,
        "ABC_123": IdentifierTok,
        "+": IdentifierTok,
        "-": IdentifierTok,
        "12a": IdentifierTok,
        "12": IntTok,
        "0": IntTok,
        "-1234": IntTok,
        "12.5": FloatTok,
        "-0.5": FloatTok,
        "let": LetTok,
        "Let": IdentifierTok,
        " if\n": IfTok,
        "iff": IdentifierTok,
    }

    for input, expected := range table {
        tokens := Tokenise(input)

        if len(tokens) != 1 {
            t.Errorf("Tokenise(\"%s\") expected to return 1 token but instead returned %d tokens", input, len(tokens))
        } else {
            tok := tokens[0]

            if tok.Type != expected {
                t.Errorf("Tokenise(\"%s\") gave token type %v but expected token type %v", input, tok.Type, expected)
            }
        }
    }
}

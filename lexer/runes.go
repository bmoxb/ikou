package lexer

func runeIsIdentChar(r rune) bool {
    return runeIsAlpha(r) || runeIsNumeral(r) || runeIsOneOf(r, "+-/*_")
}

func runeIsAlpha(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func runeIsNumeral(r rune) bool {
	return '0' <= r && r <= '9'
}

func runeIsOneOf(r rune, chars string) bool {
    allowedRunes := make(map[rune]struct{})

    for _, allowedRune := range []rune(chars) {
        allowedRunes[allowedRune] = struct{}{}
    }

    _, present := allowedRunes[r]
    return present
}

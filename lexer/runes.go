package lexer

func runeIsAlpha(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func runeIsNumeral(r rune) bool {
	return '0' <= r && r <= '9'
}

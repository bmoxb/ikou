package lexer

import (
	"fmt"
	"strings"
)

type LexicalError struct {
	pos  TokenPosition
	line string
	msg  string
}

func (l *LexicalError) Error() string {
	builder := strings.Builder{}
	for i := 0; i < int(l.pos.HorizontalPosition); i++ {
		builder.WriteRune(' ')
	}

	upper := strings.Builder{}
	upper.WriteString(builder.String())
	upper.WriteRune('▼')

	lower := strings.Builder{}
	lower.WriteString(builder.String())
	lower.WriteRune('▲')

	return fmt.Sprintf("| %s\n| %s\n| %s\n\nLexical error at %v - %s", upper.String(), l.line, lower.String(), l.pos, l.msg)
}

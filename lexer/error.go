package lexer

import (
	"fmt"

	"github.com/bmoxb/ikou/errors"
	"github.com/bmoxb/ikou/tokens"
)

type LexicalError struct {
	pos  tokens.Position
	line string
	msg  string
}

func (l *LexicalError) Error() string {
	return fmt.Sprintf("%s\n\nLexical error at %v - %s", errors.BuildRelevantLineString(l.line, l.pos), l.pos, l.msg)
}

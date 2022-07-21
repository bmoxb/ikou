package parser

import (
	"fmt"

	"github.com/bmoxb/ikou/errors"
	"github.com/bmoxb/ikou/tokens"
)

type ParsingError struct {
	pos  tokens.Position
	line string
	msg  string
}

func (l *ParsingError) Error() string {
	return fmt.Sprintf("%s\n\nParsing error at %v - %s", errors.BuildRelevantLineString(l.line, l.pos), l.pos, l.msg)
}

package errors

import (
	"fmt"
	"strings"

	"github.com/bmoxb/ikou/tokens"
)

func BuildRelevantLineString(line string, pos tokens.Position) string {
	builder := strings.Builder{}
	for i := 0; i < int(pos.HorizontalPosition); i++ {
		builder.WriteRune(' ')
	}

	upper := strings.Builder{}
	upper.WriteString(builder.String())
	upper.WriteRune('▼')

	lower := strings.Builder{}
	lower.WriteString(builder.String())
	lower.WriteRune('▲')

	return fmt.Sprintf("| %s\n| %s\n| %s", upper.String(), line, lower.String())
}

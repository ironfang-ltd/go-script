package lexer

import (
	"fmt"
	"strings"
)

type TokenError struct {
	Message string
	Source  string
	Column  int
	Line    int
}

func (e *TokenError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s on line %d at column %d: \n", e.Message, e.Line, e.Column))

	lines := strings.Split(e.Source, "\n")
	line := lines[e.Line-1]

	// Replace tabs with 4 spaces
	tabs := strings.Count(line, "\t")
	line = strings.ReplaceAll(line, "\t", "    ")

	sb.WriteString(fmt.Sprintf("%s\n", line))
	sb.WriteString(strings.Repeat("-", e.Column-1+(3*tabs)))
	sb.WriteString("^")

	return sb.String()
}

func NewTokenError(message, source string, line, column int) *TokenError {
	return &TokenError{
		Message: message,
		Source:  source,
		Line:    line,
		Column:  column,
	}
}

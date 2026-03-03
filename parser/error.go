package parser

import (
	"fmt"
	"strings"

	"github.com/ironfang-ltd/go-script/lexer"
)

type ParseError struct {
	Message string
	Source  string
	Token   lexer.Token
}

func (e *ParseError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("error: %s\n", e.Message))
	sb.WriteString(fmt.Sprintf(" --> line %d, column %d\n", e.Token.Line, e.Token.Column))

	lines := strings.Split(e.Source, "\n")

	// Bounds check
	if e.Token.Line < 1 || e.Token.Line > len(lines) {
		return sb.String()
	}

	line := lines[e.Token.Line-1]

	// Replace tabs with 4 spaces and adjust column
	tabs := strings.Count(line[:min(e.Token.Column-1, len(line))], "\t")
	line = strings.ReplaceAll(line, "\t", "    ")

	lineNumStr := fmt.Sprintf("%d", e.Token.Line)
	padding := strings.Repeat(" ", len(lineNumStr))

	sb.WriteString(fmt.Sprintf("%s |\n", padding))
	sb.WriteString(fmt.Sprintf("%s | %s\n", lineNumStr, line))

	col := e.Token.Column - 1 + (3 * tabs)
	sb.WriteString(fmt.Sprintf("%s | %s^", padding, strings.Repeat(" ", col)))

	return sb.String()
}

func NewParseError(message, source string, token lexer.Token) *ParseError {
	return &ParseError{
		Message: message,
		Source:  source,
		Token:   token,
	}
}

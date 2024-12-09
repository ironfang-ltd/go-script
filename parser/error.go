package parser

import (
	"fmt"
	"strings"

	"github.com/ironfang-ltd/ironscript/lexer"
)

type ParseError struct {
	Message string
	Source  string
	Token   *lexer.Token
}

func (e *ParseError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s on line %d, column %d: \n", e.Message, e.Token.Line, e.Token.Column))

	lines := strings.Split(e.Source, "\n")
	line := lines[e.Token.Line-1]

	// Replace tabs with 4 spaces
	tabs := strings.Count(line, "\t")
	line = strings.ReplaceAll(line, "\t", "    ")

	sb.WriteString(fmt.Sprintf("%s\n", line))
	sb.WriteString(strings.Repeat("-", e.Token.Column-1+(3*tabs)))
	sb.WriteString("^")

	return sb.String()
}

func NewParseError(message, source string, token *lexer.Token) *ParseError {
	return &ParseError{
		Message: message,
		Source:  source,
		Token:   token,
	}
}

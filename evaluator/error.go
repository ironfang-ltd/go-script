package evaluator

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	Message string
	Source  string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("error: %s\n", e.Message))
	sb.WriteString(fmt.Sprintf(" --> line %d, column %d\n", e.Line, e.Column))

	lines := strings.Split(e.Source, "\n")

	// Bounds check
	if e.Line < 1 || e.Line > len(lines) {
		return sb.String()
	}

	line := lines[e.Line-1]

	// Replace tabs with 4 spaces and adjust column
	tabs := strings.Count(line[:min(e.Column-1, len(line))], "\t")
	line = strings.ReplaceAll(line, "\t", "    ")

	lineNumStr := fmt.Sprintf("%d", e.Line)
	padding := strings.Repeat(" ", len(lineNumStr))

	sb.WriteString(fmt.Sprintf("%s |\n", padding))
	sb.WriteString(fmt.Sprintf("%s | %s\n", lineNumStr, line))

	col := e.Column - 1 + (3 * tabs)
	sb.WriteString(fmt.Sprintf("%s | %s^", padding, strings.Repeat(" ", col)))

	return sb.String()
}

func NewRuntimeError(message, source string, line, column int) *RuntimeError {
	return &RuntimeError{
		Message: message,
		Source:  source,
		Line:    line,
		Column:  column,
	}
}

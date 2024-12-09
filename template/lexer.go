package template

import (
	"github.com/ironfang-ltd/ironscript/lexer"
	"strings"
)

const (
	ForEach   = "foreach"
	If        = "if"
	Else      = "else"
	End       = "end"
	CodeStart = "{%"
	CodeEnd   = "%}"
)

type Lexer struct {
	source   string
	position int
	line     int
	col      int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:   source,
		position: 0,
		line:     1,
		col:      1,
	}
}

func (l *Lexer) Read() (lexer.Token, error) {

	if l.position >= len(l.source) {
		return lexer.NewToken(lexer.EndOfFile, "", l.position, l.line, l.col), nil
	}

	if strings.HasPrefix(l.source[l.position:], CodeStart) {
		return l.parseCode()
	}

	return l.parseText()
}

func (l *Lexer) parseCode() (lexer.Token, error) {
	start := l.position + 2
	end := start

	for {
		if l.position >= len(l.source) {
			return lexer.NewToken(lexer.Code, l.source[start:end], start, l.line, l.col), nil
		}

		if strings.HasPrefix(l.source[l.position:], CodeEnd) {
			end = l.position
			l.position += len(CodeEnd)
			break
		}

		if l.source[l.position] == '\n' {
			l.col = 1
			l.line++
		}

		l.position++
	}

	return lexer.NewToken(lexer.Code, l.source[start:end], start, l.line, l.col), nil
}

func (l *Lexer) parseText() (lexer.Token, error) {
	start := l.position

	for {
		if l.position >= len(l.source) {
			break
		}

		if strings.HasPrefix(l.source[l.position:], CodeStart) {
			break
		}

		if l.source[l.position] == '\n' {
			l.col = 1
			l.line++
		}

		l.position++
	}

	return lexer.NewToken(lexer.Text, l.source[start:l.position], start, l.line, l.col), nil
}

package lexer

import (
	"fmt"
	"strings"
)

type Mode int

const (
	ModeScript Mode = iota
	ModeTemplate
)

const (
	ScriptStartToken = "{%"
)

var keywords = map[string]TokenType{
	"as":      As,
	"let":     Let,
	"fn":      Function,
	"return":  Return,
	"true":    True,
	"false":   False,
	"if":      If,
	"else":    Else,
	"foreach": Foreach,
}

type Lexer struct {
	source        string
	position      int
	line          int
	col           int
	mode          Mode
	parseTemplate bool
}

func NewTemplate(source string) *Lexer {
	return &Lexer{
		source:        source,
		position:      0,
		line:          1,
		col:           1,
		mode:          ModeTemplate,
		parseTemplate: true,
	}
}

func NewScript(source string) *Lexer {
	return &Lexer{
		source:        source,
		position:      0,
		line:          1,
		col:           1,
		mode:          ModeScript,
		parseTemplate: false,
	}
}

func (l *Lexer) Read() (Token, error) {

	if l.position >= len(l.source) {
		return NewToken(EndOfFile, "", l.position, l.line, l.col), nil
	}

	if l.mode == ModeTemplate {
		if strings.HasPrefix(l.source[l.position:], ScriptStartToken) {
			start := l.position
			col := l.col
			l.position += len(ScriptStartToken)
			l.col += len(ScriptStartToken)
			l.mode = ModeScript
			return NewToken(ScriptStart, ScriptStartToken, start, l.line, col), nil
		}

		return l.readTemplate()
	} else {
		return l.readScript()
	}
}

func (l *Lexer) readTemplate() (Token, error) {
	start := l.position
	col := l.col
	line := l.line

	for {
		if l.position >= len(l.source) {
			break
		}

		if strings.HasPrefix(l.source[l.position:], ScriptStartToken) {
			break
		}

		if l.source[l.position] == '\n' {
			l.col = 0
			l.line++
		}

		l.col++
		l.position++
	}

	return NewToken(Text, l.source[start:l.position], start, line, col), nil
}

func (l *Lexer) readScript() (Token, error) {
	for {
		l.consumeWhitespace()

		if l.position >= len(l.source) {
			return NewToken(EndOfFile, "", l.position, l.line, l.col), nil
		}

		if l.source[l.position] == '\n' {
			l.position++
			l.line++
			l.col = 1
			continue
		}

		if l.parseTemplate {
			if token, ok := l.trySequence("%}", ScriptEnd); ok {
				l.mode = ModeTemplate
				return token, nil
			}
		}

		if token, ok := l.trySingle('(', LeftParen); ok {
			return token, nil
		}

		if token, ok := l.trySingle(')', RightParen); ok {
			return token, nil
		}

		if token, ok := l.trySingle('{', LeftBrace); ok {
			return token, nil
		}

		if token, ok := l.trySingle('}', RightBrace); ok {
			return token, nil
		}

		if token, ok := l.trySingle('[', LeftBracket); ok {
			return token, nil
		}

		if token, ok := l.trySingle(']', RightBracket); ok {
			return token, nil
		}

		if token, ok, err := l.tryString('"', String); ok || err != nil {
			return token, err
		}

		if token, ok := l.trySingle('.', Dot); ok {
			return token, nil
		}

		if token, ok := l.trySingle(',', Comma); ok {
			return token, nil
		}

		if token, ok := l.trySingle(':', Colon); ok {
			return token, nil
		}

		if token, ok := l.trySingle(';', Semicolon); ok {
			return token, nil
		}

		if token, ok := l.trySingle('+', Plus); ok {
			return token, nil
		}

		if token, ok := l.trySingle('-', Minus); ok {
			return token, nil
		}

		if token, ok := l.trySingle('%', Modulo); ok {
			return token, nil
		}

		if token, ok := l.trySingle('*', Asterisk); ok {
			return token, nil
		}

		if token, ok := l.trySingle('/', Slash); ok {
			return token, nil
		}

		if token, ok := l.trySequence("==", Equals); ok {
			return token, nil
		}

		if token, ok := l.trySequence("!=", NotEqual); ok {
			return token, nil
		}

		if token, ok := l.trySingle('=', Equal); ok {
			return token, nil
		}

		if token, ok := l.trySingle('!', Bang); ok {
			return token, nil
		}

		if token, ok := l.trySequence("<=", LessOrEqual); ok {
			return token, nil
		}

		if token, ok := l.trySequence(">=", GreaterOrEqual); ok {
			return token, nil
		}

		if token, ok := l.trySingle('<', LessThan); ok {
			return token, nil
		}

		if token, ok := l.trySingle('>', GreaterThan); ok {
			return token, nil
		}

		if token, ok := l.tryNumber(); ok {
			return token, nil
		}

		if token, ok := l.tryIdentifierOrKeyword(); ok {
			return token, nil
		}

		break
	}

	return Token{}, NewTokenError(
		fmt.Sprintf("unexpected character '%c'", l.source[l.position]),
		l.source, l.line, l.col)
}

func (l *Lexer) GetSource() string {
	return l.source
}

func (l *Lexer) consumeWhitespace() {
	for l.position < len(l.source) && (l.source[l.position] == ' ' || l.source[l.position] == '\r' || l.source[l.position] == '\t') {
		l.position++
		l.col++
	}
}

func (l *Lexer) trySingle(c byte, tokenType TokenType) (Token, bool) {

	pos := l.position
	col := l.col
	line := l.line

	if l.source[l.position] == c {
		l.position++
		l.col++
		return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true
	}

	return TokenNone, false
}

func (l *Lexer) trySequence(seq string, tokenType TokenType) (Token, bool) {

	pos := l.position
	col := l.col
	line := l.line

	if strings.HasPrefix(l.source[l.position:], seq) {
		l.position += len(seq)
		l.col += len(seq)
		return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true
	}

	return TokenNone, false
}

func (l *Lexer) tryString(quote byte, tokenType TokenType) (Token, bool, error) {

	pos := l.position
	col := l.col
	line := l.line

	if l.source[l.position] != quote {
		return TokenNone, false, nil
	}

	l.position++
	l.col++

	for {
		if l.position >= len(l.source) {
			return TokenNone, false, NewTokenError("unterminated string literal", l.source, line, col)
		}

		if l.source[l.position] == '\\' {
			// Skip escape sequence (backslash + next character)
			l.position++
			l.col++
			if l.position < len(l.source) {
				l.position++
				l.col++
			}
			continue
		}

		if l.source[l.position] == quote {
			l.position++
			l.col++
			return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true, nil
		}

		if l.source[l.position] == '\n' {
			return TokenNone, false, NewTokenError("unterminated string literal", l.source, line, col)
		}

		l.position++
		l.col++
	}
}

func (l *Lexer) tryNumber() (Token, bool) {

	ok := false
	tokenType := Integer
	hasDot := false

	pos := l.position
	col := l.col
	line := l.line

	for {
		if l.position >= len(l.source) {
			break
		}

		if l.source[l.position] >= '0' && l.source[l.position] <= '9' {
			l.position++
			l.col++
			ok = true
			continue
		}

		if l.source[l.position] == '.' && !hasDot {
			// Only consume dot if followed by a digit
			if l.position+1 < len(l.source) && l.source[l.position+1] >= '0' && l.source[l.position+1] <= '9' {
				l.position++
				l.col++
				tokenType = Float
				hasDot = true
				continue
			}
			break
		}

		break
	}

	if !ok {
		return Token{}, false
	}

	return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true
}

func (l *Lexer) tryIdentifierOrKeyword() (Token, bool) {

	pos := l.position
	col := l.col
	line := l.line

	c := l.source[l.position]
	if !(c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
		return TokenNone, false
	}

	for l.position < len(l.source) {
		c = l.source[l.position]
		if c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			l.position++
			l.col++
			continue
		}
		break
	}

	word := l.source[pos:l.position]
	if tokenType, ok := keywords[word]; ok {
		return NewToken(tokenType, word, pos, line, col), true
	}

	return NewToken(Identifier, word, pos, line, col), true
}

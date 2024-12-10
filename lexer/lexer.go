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
	ScriptEndToken   = "%}"
)

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

		pos := l.position
		col := l.col
		line := l.line

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

		if token, ok := l.trySequence("as", As); ok {
			return token, nil
		}

		if token, ok := l.trySequence("let", Let); ok {
			return token, nil
		}

		if token, ok := l.trySequence("fn", Function); ok {
			return token, nil
		}

		if token, ok := l.trySequence("return", Return); ok {
			return token, nil
		}

		if token, ok := l.trySequence("true", True); ok {
			return token, nil
		}

		if token, ok := l.trySequence("false", False); ok {
			return token, nil
		}

		if token, ok := l.trySequence("if", If); ok {
			return token, nil
		}

		if token, ok := l.trySequence("else", Else); ok {
			return token, nil
		}

		if token, ok := l.trySequence("foreach", Foreach); ok {
			return token, nil
		}

		if token, ok := l.tryNumber(); ok {
			return token, nil
		}

		if l.tryIdentifier() {
			return NewToken(Identifier, l.source[pos:l.position], pos, line, col), nil
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
			return TokenNone, false, NewTokenError("unexpected end of file when parsing string", l.source, line, col)
		}

		if l.source[l.position] == quote {

			// check for escaped quote
			if l.position > 0 && l.source[l.position-1] == '\\' {
				l.position++
				l.col++
				continue
			}

			l.position++
			l.col++
			return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true, nil
		}

		if l.source[l.position] == '\n' {
			return TokenNone, false, nil
		}

		l.position++
		l.col++
	}
}

func (l *Lexer) tryNumber() (Token, bool) {

	ok := false
	tokenType := Integer

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

		if l.source[l.position] == '.' {
			l.position++
			l.col++
			tokenType = Float
			ok = true
			continue
		}

		break
	}

	if !ok {
		return Token{}, false
	}

	return NewToken(tokenType, l.source[pos:l.position], pos, line, col), true
}

func (l *Lexer) tryIdentifier() bool {

	ok := false

	for {
		if l.position >= len(l.source) {
			break
		}

		if l.source[l.position] == '_' || (l.source[l.position] >= 'a' && l.source[l.position] <= 'z') || (l.source[l.position] >= 'A' && l.source[l.position] <= 'Z') {
			l.position++
			l.col++
			ok = true
			continue
		}

		break
	}

	return ok
}

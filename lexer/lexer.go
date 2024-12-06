package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	source   string
	position int
	line     int
	col      int
}

func New(source string) *Lexer {
	return &Lexer{
		source:   source,
		position: 0,
		line:     1,
		col:      1,
	}
}

func (l *Lexer) Read() (Token, error) {
	for {
		l.consumeWhitespace()

		if l.position >= len(l.source) {
			return NewToken(EndOfFile, "", l.line, l.col), nil
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

		if token, ok := l.trySingle('=', Equal); ok {
			return token, nil
		}

		if token, ok := l.trySequence("!=", NotEqual); ok {
			return token, nil
		}

		if token, ok := l.trySingle('!', Bang); ok {
			return token, nil
		}

		if token, ok := l.trySequence("==", Equals); ok {
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

		if token, ok := l.trySequence("for", For); ok {
			return token, nil
		}

		if token, ok := l.tryNumber(); ok {
			return token, nil
		}

		if l.tryIdentifier() {
			return NewToken(Identifier, l.source[pos:l.position], line, col), nil
		}

		break
	}

	return Token{}, fmt.Errorf("unexpected character '%c' at line %d, col %d", l.source[l.position], l.line, l.col)
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
		return NewToken(tokenType, l.source[pos:l.position], col, line), true
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
		return NewToken(tokenType, l.source[pos:l.position], col, line), true
	}

	return TokenNone, false
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
		}

		break
	}

	if !ok {
		return Token{}, false
	}

	return NewToken(tokenType, l.source[pos:l.position], line, col), true
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

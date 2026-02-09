package lexer

import (
	"fmt"
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
		if l.source[l.position] == '{' && l.position+1 < len(l.source) && l.source[l.position+1] == '%' {
			start := l.position
			col := l.col
			l.position += 2
			l.col += 2
			l.mode = ModeScript
			return NewToken(ScriptStart, l.source[start:l.position], start, l.line, col), nil
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

	for l.position < len(l.source) {
		if l.source[l.position] == '{' && l.position+1 < len(l.source) && l.source[l.position+1] == '%' {
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

		ch := l.source[l.position]

		if ch == '\n' {
			l.position++
			l.line++
			l.col = 1
			continue
		}

		pos := l.position
		col := l.col
		line := l.line

		switch ch {
		case '(':
			l.position++
			l.col++
			return NewToken(LeftParen, l.source[pos:l.position], pos, line, col), nil
		case ')':
			l.position++
			l.col++
			return NewToken(RightParen, l.source[pos:l.position], pos, line, col), nil
		case '{':
			l.position++
			l.col++
			return NewToken(LeftBrace, l.source[pos:l.position], pos, line, col), nil
		case '}':
			l.position++
			l.col++
			return NewToken(RightBrace, l.source[pos:l.position], pos, line, col), nil
		case '[':
			l.position++
			l.col++
			return NewToken(LeftBracket, l.source[pos:l.position], pos, line, col), nil
		case ']':
			l.position++
			l.col++
			return NewToken(RightBracket, l.source[pos:l.position], pos, line, col), nil
		case '.':
			l.position++
			l.col++
			return NewToken(Dot, l.source[pos:l.position], pos, line, col), nil
		case ',':
			l.position++
			l.col++
			return NewToken(Comma, l.source[pos:l.position], pos, line, col), nil
		case ':':
			l.position++
			l.col++
			return NewToken(Colon, l.source[pos:l.position], pos, line, col), nil
		case ';':
			l.position++
			l.col++
			return NewToken(Semicolon, l.source[pos:l.position], pos, line, col), nil
		case '+':
			l.position++
			l.col++
			return NewToken(Plus, l.source[pos:l.position], pos, line, col), nil
		case '-':
			l.position++
			l.col++
			return NewToken(Minus, l.source[pos:l.position], pos, line, col), nil
		case '*':
			l.position++
			l.col++
			return NewToken(Asterisk, l.source[pos:l.position], pos, line, col), nil
		case '/':
			l.position++
			l.col++
			return NewToken(Slash, l.source[pos:l.position], pos, line, col), nil
		case '%':
			if l.parseTemplate && l.position+1 < len(l.source) && l.source[l.position+1] == '}' {
				l.position += 2
				l.col += 2
				l.mode = ModeTemplate
				return NewToken(ScriptEnd, l.source[pos:l.position], pos, line, col), nil
			}
			l.position++
			l.col++
			return NewToken(Modulo, l.source[pos:l.position], pos, line, col), nil
		case '=':
			if l.position+1 < len(l.source) && l.source[l.position+1] == '=' {
				l.position += 2
				l.col += 2
				return NewToken(Equals, l.source[pos:l.position], pos, line, col), nil
			}
			l.position++
			l.col++
			return NewToken(Equal, l.source[pos:l.position], pos, line, col), nil
		case '!':
			if l.position+1 < len(l.source) && l.source[l.position+1] == '=' {
				l.position += 2
				l.col += 2
				return NewToken(NotEqual, l.source[pos:l.position], pos, line, col), nil
			}
			l.position++
			l.col++
			return NewToken(Bang, l.source[pos:l.position], pos, line, col), nil
		case '<':
			if l.position+1 < len(l.source) && l.source[l.position+1] == '=' {
				l.position += 2
				l.col += 2
				return NewToken(LessOrEqual, l.source[pos:l.position], pos, line, col), nil
			}
			l.position++
			l.col++
			return NewToken(LessThan, l.source[pos:l.position], pos, line, col), nil
		case '>':
			if l.position+1 < len(l.source) && l.source[l.position+1] == '=' {
				l.position += 2
				l.col += 2
				return NewToken(GreaterOrEqual, l.source[pos:l.position], pos, line, col), nil
			}
			l.position++
			l.col++
			return NewToken(GreaterThan, l.source[pos:l.position], pos, line, col), nil
		case '"':
			token, _, err := l.tryString('"', String)
			return token, err
		default:
			if ch >= '0' && ch <= '9' {
				token, _ := l.tryNumber()
				return token, nil
			}
			if ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				token, _ := l.tryIdentifierOrKeyword()
				return token, nil
			}
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

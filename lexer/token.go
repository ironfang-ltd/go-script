package lexer

type TokenType string

const (
	None           TokenType = "NONE"
	EndOfFile      TokenType = "EOF"
	Identifier     TokenType = "IDENTIFIER"
	String         TokenType = "STRING"
	LeftParen      TokenType = "LEFT_PAREN"
	RightParen     TokenType = "RIGHT_PAREN"
	LeftBrace      TokenType = "LEFT_BRACE"
	RightBrace     TokenType = "RIGHT_BRACE"
	LeftBracket    TokenType = "LEFT_BRACKET"
	RightBracket   TokenType = "RIGHT_BRACKET"
	Comma          TokenType = "COMMA"
	Colon          TokenType = "COLON"
	Semicolon      TokenType = "SEMICOLON"
	Integer        TokenType = "INTEGER"
	Float          TokenType = "FLOAT"
	Plus           TokenType = "PLUS"
	Minus          TokenType = "MINUS"
	Modulo         TokenType = "MODULO"
	Asterisk       TokenType = "ASTERISK"
	Dot            TokenType = "DOT"
	Slash          TokenType = "SLASH"
	Equal          TokenType = "EQUAL"
	Equals         TokenType = "EQUALS"
	NotEqual       TokenType = "NOT_EQUAL"
	LessThan       TokenType = "LESS_THAN"
	GreaterThan    TokenType = "GREATER_THAN"
	LessOrEqual    TokenType = "LESS_OR_EQUAL"
	GreaterOrEqual TokenType = "GREATER_OR_EQUAL"
	Function       TokenType = "FUNCTION"
	Let            TokenType = "LET"
	Return         TokenType = "RETURN"
	True           TokenType = "TRUE"
	False          TokenType = "FALSE"
	If             TokenType = "IF"
	Else           TokenType = "ELSE"
	Bang           TokenType = "BANG"
	Foreach        TokenType = "FOREACH"
	As             TokenType = "AS"
	Text           TokenType = "TEXT"
	ScriptStart    TokenType = "SCRIPT_START"
	ScriptEnd      TokenType = "SCRIPT_END"
)

var TokenNone = NewToken(None, "", 0, 0, 0)

type Token struct {
	Type     TokenType
	Source   string
	Position int
	Line     int
	Column   int
}

func NewToken(typ TokenType, source string, position, line, column int) Token {
	return Token{
		Type:     typ,
		Source:   source,
		Position: position,
		Line:     line,
		Column:   column,
	}
}

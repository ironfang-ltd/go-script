package parser

import (
	"errors"
	"fmt"
	"github.com/ironfang-ltd/go-script/lexer"
	"strconv"
	"strings"
)

var Precedences = []lexer.TokenType{
	lexer.None,
	lexer.Equals,
	lexer.NotEqual,
	lexer.LessThan,
	lexer.GreaterThan,
	lexer.Plus,
	lexer.Minus,
	lexer.Slash,
	lexer.Asterisk,
	lexer.Dot,
	lexer.Modulo,
	//
	lexer.LeftParen,
	lexer.LeftBracket,
}

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: []Statement{},
	}
}

func (p *Program) Debug() string {
	str := ""
	for _, s := range p.Statements {
		str += s.Debug() + "\n"
	}
	return str
}

type Parser struct {
	l       *lexer.Lexer
	prev    *lexer.Token
	current *lexer.Token
	next    *lexer.Token
	errors  []error
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		l: l,
	}
}

func (p *Parser) Parse() (*Program, error) {

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	t := NewProgram()

	for {

		if p.current.Type == lexer.EndOfFile {
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		if statement != nil {
			t.Statements = append(t.Statements, statement)
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}
	}

	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
	}

	return t, nil
}

func (p *Parser) parseStatement() (Statement, error) {

	switch p.current.Type {
	case lexer.ScriptStart:
		return nil, nil
	case lexer.ScriptEnd:
		return nil, nil
	case lexer.Text:
		return p.parseTextStatement()
	case lexer.Let:
		return p.parseLetStatement()
	case lexer.Foreach:
		return p.parseForeachExpression()
	case lexer.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseTextStatement() (*PrintStatement, error) {
	statement := &PrintStatement{
		Value: p.current.Source,
	}

	if p.next.Type == lexer.ScriptStart {
		lastNewLine := strings.LastIndexByte(statement.Value, '\n')
		if lastNewLine != -1 {
			// Remove the last new line and whitespace after it
			// if the last line is empty
			if strings.TrimSpace(statement.Value[lastNewLine:]) == "" {
				statement.Value = statement.Value[:lastNewLine]
			}
		}
	}

	return statement, nil
}

func (p *Parser) parseLetStatement() (*LetStatement, error) {

	statement := &LetStatement{
		Token: *p.current,
	}

	peek, err := p.tryPeek(lexer.Identifier)
	if !peek || err != nil {
		return nil, err
	}

	statement.Name = &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	}

	peek, err = p.tryPeek(lexer.Equal)
	if !peek || err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpression(0)
	if value == nil || err != nil {
		return nil, err
	}

	statement.Value = value

	if p.next.Type != lexer.Semicolon && p.next.Type != lexer.ScriptEnd && p.next.Type != lexer.EndOfFile {

		p.errors = append(p.errors,
			NewParseError(
				fmt.Sprintf("expected %s, %s or %s, got %s", lexer.Semicolon, lexer.ScriptEnd, lexer.EndOfFile, p.current.Type),
				p.l.GetSource(), p.current))

		return nil, nil
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func (p *Parser) parseReturnStatement() (*ReturnStatement, error) {
	statement := &ReturnStatement{
		Token: *p.current,
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpression(0)
	if value == nil || err != nil {
		return nil, nil
	}

	statement.Value = value

	peek, err := p.tryPeek(lexer.Semicolon)
	if !peek || err != nil {
		return nil, err
	}

	return statement, nil
}

func (p *Parser) parseExpressionStatement() (*ExpressionStatement, error) {
	expression, err := p.parseExpression(0)
	if expression == nil || err != nil {
		return nil, err
	}

	statement := &ExpressionStatement{
		Expression: expression,
	}

	if p.next.Type == lexer.Equal {

		switch expression.(type) {
		case *IndexExpression, *PropertyExpression, *Identifier:
			err = p.nextToken()
			if err != nil {
				return nil, err
			}

			err = p.nextToken()
			if err != nil {
				return nil, err
			}

			right, err := p.parseExpression(0)
			if right == nil || err != nil {
				return nil, err
			}

			// This is an assignment statement
			assignment := &AssignmentExpression{
				Token: *p.current,
				Left:  expression,
				Right: right,
			}

			statement.Expression = assignment

			err = p.nextToken()
			if err != nil {
				return nil, err
			}

			return statement, nil
		}

	}

	if _, ok := expression.(*FunctionLiteral); ok {
		return statement, nil
	}

	if _, ok := expression.(*IfExpression); ok {
		return statement, nil
	}

	if p.next.Type != lexer.Semicolon && p.next.Type != lexer.ScriptEnd && p.next.Type != lexer.EndOfFile {

		p.errors = append(p.errors,
			NewParseError(
				fmt.Sprintf("expected %s or %s, got %s", lexer.Semicolon, lexer.ScriptEnd, p.next.Type),
				p.l.GetSource(), p.next))

		return nil, nil
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	return statement, nil
}

func (p *Parser) parseExpression(precedence int) (Expression, error) {

	leftExpression, err := p.parsePrefixExpression()
	if leftExpression == nil || err != nil {
		return nil, err
	}

	for {
		if (p.next.Type == lexer.Semicolon || p.next.Type == lexer.ScriptEnd) || precedence >= p.peekPrecedence() {
			return leftExpression, nil
		}

		switch p.next.Type {
		case lexer.Plus:
			fallthrough
		case lexer.Minus:
			fallthrough
		case lexer.Asterisk:
			fallthrough
		case lexer.Slash:
			fallthrough
		case lexer.Equals:
			fallthrough
		case lexer.NotEqual:
			fallthrough
		case lexer.GreaterThan:
			fallthrough
		case lexer.LessThan:
			err := p.nextToken()
			if err != nil {
				return nil, err
			}

			leftExpression, err = p.parseInfixExpression(leftExpression)
			if err != nil {
				return nil, err
			}
		case lexer.Dot:
			err := p.nextToken()
			if err != nil {
				return nil, err
			}

			leftExpression, err = p.parseAccessExpression(leftExpression)
			if err != nil {
				return nil, err
			}
		case lexer.LeftParen:
			err := p.nextToken()
			if err != nil {
				return nil, err
			}

			leftExpression, err = p.parseCallExpression(leftExpression)
			if err != nil {
				return nil, err
			}
		case lexer.LeftBracket:
			err := p.nextToken()
			if err != nil {
				return nil, err
			}

			leftExpression, err = p.parseIndexExpression(leftExpression)
			if err != nil {
				return nil, err
			}
		default:
			return leftExpression, nil
		}

		if leftExpression == nil {
			break
		}
	}

	return leftExpression, nil
}

func (p *Parser) parsePrefixExpression() (Expression, error) {
	switch p.current.Type {
	case lexer.Identifier:
		return p.parseIdentifier()
	case lexer.Integer:
		return p.parseInteger()
	case lexer.String:
		return p.parseString()
	case lexer.True:
		return p.parseBoolean()
	case lexer.False:
		return p.parseBoolean()
	case lexer.Bang:
		fallthrough
	case lexer.Minus:
		expression := &PrefixExpression{
			Token:    *p.current,
			Operator: p.current.Source,
		}

		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		right, err := p.parseExpression(5) // Plus
		if right == nil || err != nil {
			return nil, err
		}

		expression.Right = right

		return expression, nil
	case lexer.LeftParen:
		return p.parseGroupedExpression()
	case lexer.If:
		return p.parseIfExpression()
	case lexer.Function:
		return p.parseFunctionLiteral()
	case lexer.LeftBracket:
		return p.parseArray()
	case lexer.LeftBrace:
		return p.parseHashLiteral()
	default:
		p.errors = append(p.errors,
			NewParseError(fmt.Sprintf("unexpected token %s", p.current.Type), p.l.GetSource(), p.current))
		return nil, nil
	}
}

func (p *Parser) parseInfixExpression(left Expression) (Expression, error) {

	infix := &InfixExpression{
		Token: *p.current,
		Left:  left,
	}

	precedence := p.currentPrecedence()

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	right, err := p.parseExpression(precedence)
	if right == nil || err != nil {
		return nil, err
	}

	infix.Right = right

	return infix, nil
}

func (p *Parser) parseAccessExpression(left Expression) (Expression, error) {

	expression := &PropertyExpression{
		Token: *p.current,
		Left:  left,
	}

	peek, err := p.tryPeek(lexer.Identifier)
	if !peek || err != nil {
		return nil, err
	}

	exp, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	expression.Property = exp

	return expression, nil
}

func (p *Parser) parseCallExpression(left Expression) (Expression, error) {
	args, err := p.parseExpressionList(lexer.RightParen)
	if err != nil {
		return nil, err
	}

	expression := &CallExpression{
		Token:    *p.current,
		Function: left,
		Args:     args,
	}

	return expression, nil
}

func (p *Parser) parseIndexExpression(left Expression) (Expression, error) {

	expression := &IndexExpression{
		Token: *p.current,
		Left:  left,
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	exp, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	expression.Index = exp

	peek, err := p.tryPeek(lexer.RightBracket)
	if !peek || err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) currentPrecedence() int {

	for i := 0; i < len(Precedences); i++ {
		if Precedences[i] == p.current.Type {
			return i
		}
	}

	return 0
}

func (p *Parser) peekPrecedence() int {

	for i := 0; i < len(Precedences); i++ {
		if Precedences[i] == p.next.Type {
			return i
		}
	}

	return 0
}

func (p *Parser) parseIdentifier() (*Identifier, error) {
	return &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	}, nil
}

func (p *Parser) parseInteger() (Expression, error) {
	literal := &IntegerLiteral{Token: *p.current}

	i, err := strconv.ParseInt(p.current.Source, 10, 64)
	if err != nil {
		return nil, err
	}

	literal.Value = int(i)

	return literal, nil
}

func (p *Parser) parseString() (Expression, error) {

	literal := &StringLiteral{Token: *p.current}

	v := strings.Trim(p.current.Source, "\"")

	v = strings.ReplaceAll(v, "\\n", "\n")
	v = strings.ReplaceAll(v, "\\r", "\r")
	v = strings.ReplaceAll(v, "\\t", "\t")

	literal.Value = v

	return literal, nil
}

func (p *Parser) parseBoolean() (Expression, error) {

	literal := &BooleanLiteral{Token: *p.current}

	if p.current.Type == lexer.True {
		literal.Value = true
	}

	if p.current.Type == lexer.False {
		literal.Value = false
	}

	return literal, nil
}

func (p *Parser) parseGroupedExpression() (Expression, error) {

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	expression, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	peek, err := p.tryPeek(lexer.RightParen)
	if !peek || err != nil {
		return nil, err
	}

	return expression, nil
}

func (p *Parser) parseIfExpression() (Expression, error) {

	expression := &IfExpression{
		Token: *p.current,
	}

	peek, err := p.tryPeek(lexer.LeftParen)
	if !peek || err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expression.Condition, err = p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	peek, err = p.tryPeek(lexer.RightParen)
	if !peek || err != nil {
		return nil, err
	}

	peek, err = p.tryPeek(lexer.LeftBrace)
	if !peek || err != nil {
		return nil, err
	}

	expression.Consequence = p.parseBlockStatement()

	if p.next.Type == lexer.Else {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		peek, err := p.tryPeek(lexer.LeftBrace)
		if !peek || err != nil {
			return nil, err
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression, nil
}

func (p *Parser) parseForeachExpression() (Expression, error) {

	expression := &ForeachExpression{
		Token: *p.current,
	}

	peek, err := p.tryPeek(lexer.LeftParen)
	if !peek || err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expression.Iterable, err = p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	peek, err = p.tryPeek(lexer.As)
	if !peek || err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expression.Variable, err = p.parseIdentifier()

	peek, err = p.tryPeek(lexer.RightParen)
	if !peek || err != nil {
		return nil, err
	}

	peek, err = p.tryPeek(lexer.LeftBrace)
	if !peek || err != nil {
		return nil, err
	}

	expression.Body = p.parseBlockStatement()

	return expression, nil
}

func (p *Parser) parseBlockStatement() *BlockStatement {

	block := &BlockStatement{
		Token: *p.current,
	}

	err := p.nextToken()
	if err != nil {
		return nil
	}

	for {

		if p.current.Type == lexer.RightBrace || p.current.Type == lexer.EndOfFile {
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return nil
		}

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		err = p.nextToken()
		if err != nil {
			return nil
		}
	}

	return block
}

func (p *Parser) parseFunctionParameters() ([]*Identifier, error) {
	var identifiers []*Identifier

	if p.next.Type == lexer.RightParen {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return identifiers, nil
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	identifiers = append(identifiers, &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	})

	for p.next.Type == lexer.Comma {
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		identifiers = append(identifiers, &Identifier{
			Token: *p.current,
			Value: p.current.Source,
		})
	}

	peek, err := p.tryPeek(lexer.RightParen)
	if !peek || err != nil {
		return nil, err
	}

	return identifiers, nil
}

func (p *Parser) parseFunctionLiteral() (Expression, error) {

	literal := &FunctionLiteral{
		Token: *p.current,
	}

	peek, err := p.tryPeek(lexer.Identifier)
	if !peek || err != nil {
		literal.Identifier, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}

	peek, err = p.tryPeek(lexer.LeftParen)
	if !peek || err != nil {
		return nil, err
	}

	literal.Parameters, err = p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	peek, err = p.tryPeek(lexer.LeftBrace)
	if !peek || err != nil {
		return nil, err
	}

	literal.Body = p.parseBlockStatement()

	return literal, nil
}

func (p *Parser) parseExpressionList(end lexer.TokenType) ([]Expression, error) {
	var list []Expression

	if p.next.Type == end {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return list, nil
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	exp, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}
	list = append(list, exp)

	for p.next.Type == lexer.Comma {

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		exp, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}
		list = append(list, exp)
	}

	peek, err := p.tryPeek(end)
	if !peek || err != nil {
		return nil, err
	}

	return list, nil
}

func (p *Parser) parseArray() (Expression, error) {

	exp, err := p.parseExpressionList(lexer.RightBracket)
	if err != nil {
		return nil, err
	}

	literal := &ArrayLiteral{
		Token:    *p.current,
		Elements: exp,
	}

	return literal, nil
}

func (p *Parser) parseHashLiteral() (Expression, error) {

	pairs, err := p.parseHashPairs()
	if err != nil {
		return nil, err
	}

	literal := &HashLiteral{
		Token: *p.current,
		Pairs: pairs,
	}

	return literal, nil
}

func (p *Parser) parseHashPairs() (map[Expression]Expression, error) {

	pairs := make(map[Expression]Expression)

	if p.next.Type == lexer.RightBrace {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
		return pairs, nil
	}

	for {

		// parse key, must be a string for now
		peek, err := p.tryPeek(lexer.String)
		if !peek || err != nil {
			return nil, err
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		// parse the ':' separator
		peek, err = p.tryPeek(lexer.Colon)
		if !peek || err != nil {
			return nil, err
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		// parse value
		exp, err := p.parseExpression(0)
		if err != nil {
			return nil, err
		}

		pairs[key] = exp

		if p.next.Type != lexer.Comma {
			break
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}
	}

	peek, err := p.tryPeek(lexer.RightBrace)
	if !peek || err != nil {
		return nil, err
	}

	return pairs, nil
}

func (p *Parser) tryPeek(tokenToken lexer.TokenType) (bool, error) {
	if p.next.Type == tokenToken {
		err := p.nextToken()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	p.errors = append(p.errors,
		NewParseError(
			fmt.Sprintf("expected %s, got %s", tokenToken, p.next.Type), p.l.GetSource(), p.next))

	return false, nil
}

func (p *Parser) nextToken() error {

	p.prev = p.current
	p.current = p.next

	next, err := p.l.Read()
	if err != nil {
		return err
	}

	p.next = &next
	return nil
}

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
	//
	lexer.LeftParen,
	lexer.LeftBrace,
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

	p.nextToken()
	p.nextToken()

	t := NewProgram()

	for {

		if p.current.Type == lexer.EndOfFile {
			break
		}

		statement := p.parseStatement()
		if statement != nil {
			t.Statements = append(t.Statements, statement)
		}

		p.nextToken()
	}

	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
	}

	return t, nil
}

func (p *Parser) parseStatement() Statement {

	switch p.current.Type {
	case lexer.ScriptStart:
		return nil
	case lexer.ScriptEnd:
		return nil
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

func (p *Parser) parseTextStatement() *PrintStatement {
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

	return statement
}

func (p *Parser) parseLetStatement() *LetStatement {

	statement := &LetStatement{
		Token: *p.current,
	}

	if !p.tryPeek(lexer.Identifier) {
		return nil
	}

	statement.Name = &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	}

	if !p.tryPeek(lexer.Equal) {
		return nil
	}

	p.nextToken()

	value := p.parseExpression(0)
	if value == nil {
		return nil
	}

	statement.Value = value

	if p.next.Type != lexer.Semicolon && p.next.Type != lexer.ScriptEnd {

		p.errors = append(p.errors,
			NewParseError(
				fmt.Sprintf("expected %s or %s, got %s", lexer.Semicolon, lexer.ScriptEnd, p.current.Type),
				p.l.GetSource(), p.current))

		return nil
	}

	p.nextToken()

	return statement
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	statement := &ReturnStatement{
		Token: *p.current,
	}

	p.nextToken()

	value := p.parseExpression(0)
	if value == nil {
		return nil
	}

	statement.Value = value

	if !p.tryPeek(lexer.Semicolon) {
		return nil
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	expression := p.parseExpression(0)
	if expression == nil {
		return nil
	}

	statement := &ExpressionStatement{
		Expression: expression,
	}

	if _, ok := expression.(*FunctionLiteral); ok {
		return statement
	}

	if _, ok := expression.(*IfExpression); ok {
		return statement
	}

	if p.next.Type != lexer.Semicolon && p.next.Type != lexer.ScriptEnd && p.next.Type != lexer.EndOfFile {

		p.errors = append(p.errors,
			NewParseError(
				fmt.Sprintf("expected %s or %s, got %s", lexer.Semicolon, lexer.ScriptEnd, p.next.Type),
				p.l.GetSource(), p.next))

		return nil
	}

	p.nextToken()

	return statement
}

func (p *Parser) parseExpression(precedence int) Expression {

	leftExpression := p.parsePrefixExpression()
	if leftExpression == nil {
		return nil
	}

	for {
		if (p.next.Type == lexer.Semicolon || p.next.Type == lexer.ScriptEnd) || precedence >= p.peekPrecedence() {
			return leftExpression
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
			p.nextToken()
			leftExpression = p.parseInfixExpression(leftExpression)
		case lexer.Dot:
			p.nextToken()
			leftExpression = p.parseAccessExpression(leftExpression)
		case lexer.LeftParen:
			p.nextToken()
			leftExpression = p.parseCallExpression(leftExpression)
		case lexer.LeftBrace:
			p.nextToken()
			leftExpression = p.parseIndexExpression(leftExpression)
		default:
			return leftExpression
		}

		if leftExpression == nil {
			break
		}
	}

	return leftExpression
}

func (p *Parser) parsePrefixExpression() Expression {
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

		p.nextToken()

		right := p.parseExpression(5) // Plus
		if right == nil {
			return nil
		}

		expression.Right = right

		return expression
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
		return nil
	}
}

func (p *Parser) parseInfixExpression(left Expression) Expression {

	infix := &InfixExpression{
		Token: *p.current,
		Left:  left,
	}

	precedence := p.currentPrecedence()

	p.nextToken()

	right := p.parseExpression(precedence)
	if right == nil {
		return nil
	}

	infix.Right = right

	return infix
}

func (p *Parser) parseAccessExpression(left Expression) Expression {

	expression := &PropertyExpression{
		Token: *p.current,
		Left:  left,
	}

	if !p.tryPeek(lexer.Identifier) {
		return nil
	}

	expression.Property = p.parseExpression(0)

	return expression
}

func (p *Parser) parseCallExpression(left Expression) Expression {
	args := p.parseExpressionList(lexer.RightParen)

	expression := &CallExpression{
		Token:    *p.current,
		Function: left,
		Args:     args,
	}

	return expression
}

func (p *Parser) parseIndexExpression(left Expression) Expression {

	expression := &IndexExpression{
		Token: *p.current,
		Left:  left,
	}

	p.nextToken()

	expression.Index = p.parseExpression(0)

	if !p.tryPeek(lexer.RightBrace) {
		return nil
	}

	return expression
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

func (p *Parser) parseIdentifier() *Identifier {
	return &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	}
}

func (p *Parser) parseInteger() Expression {
	literal := &IntegerLiteral{Token: *p.current}

	if i, err := strconv.ParseInt(p.current.Source, 10, 64); err == nil {
		literal.Value = int(i)
	}

	return literal
}

func (p *Parser) parseString() Expression {

	literal := &StringLiteral{Token: *p.current}

	v := strings.Trim(p.current.Source, "\"")

	v = strings.ReplaceAll(v, "\\n", "\n")
	v = strings.ReplaceAll(v, "\\r", "\r")
	v = strings.ReplaceAll(v, "\\t", "\t")

	literal.Value = v

	return literal
}

func (p *Parser) parseBoolean() Expression {

	literal := &BooleanLiteral{Token: *p.current}

	if p.current.Type == lexer.True {
		literal.Value = true
	}

	if p.current.Type == lexer.False {
		literal.Value = false
	}

	return literal
}

func (p *Parser) parseGroupedExpression() Expression {

	p.nextToken()

	expression := p.parseExpression(0)

	if !p.tryPeek(lexer.RightParen) {
		return nil
	}

	return expression
}

func (p *Parser) parseIfExpression() Expression {

	expression := &IfExpression{
		Token: *p.current,
	}

	if !p.tryPeek(lexer.LeftParen) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(0)

	if !p.tryPeek(lexer.RightParen) {
		return nil
	}

	if !p.tryPeek(lexer.LeftBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.next.Type == lexer.Else {
		p.nextToken()

		if !p.tryPeek(lexer.LeftBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseForeachExpression() Expression {

	expression := &ForeachExpression{
		Token: *p.current,
	}

	if !p.tryPeek(lexer.LeftParen) {
		return nil
	}

	p.nextToken()

	expression.Iterable = p.parseExpression(0)

	if !p.tryPeek(lexer.As) {
		return nil
	}

	p.nextToken()

	expression.Variable = p.parseIdentifier()

	if !p.tryPeek(lexer.RightParen) {
		return nil
	}

	if !p.tryPeek(lexer.LeftBrace) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseBlockStatement() *BlockStatement {

	block := &BlockStatement{
		Token: *p.current,
	}

	p.nextToken()

	for {

		if p.current.Type == lexer.RightBrace || p.current.Type == lexer.EndOfFile {
			break
		}

		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	var identifiers []*Identifier

	if p.next.Type == lexer.RightParen {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	identifiers = append(identifiers, &Identifier{
		Token: *p.current,
		Value: p.current.Source,
	})

	for p.next.Type == lexer.Comma {
		p.nextToken()
		p.nextToken()

		identifiers = append(identifiers, &Identifier{
			Token: *p.current,
			Value: p.current.Source,
		})
	}

	if !p.tryPeek(lexer.RightParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionLiteral() Expression {

	literal := &FunctionLiteral{
		Token: *p.current,
	}

	if p.tryPeek(lexer.Identifier) {
		literal.Identifier = p.parseIdentifier()
	}

	if !p.tryPeek(lexer.LeftParen) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.tryPeek(lexer.LeftBrace) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []Expression {
	var list []Expression

	if p.next.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(0))

	for p.next.Type == lexer.Comma {
		p.nextToken()
		p.nextToken()

		list = append(list, p.parseExpression(0))
	}

	p.tryPeek(end)

	return list
}

func (p *Parser) parseArray() Expression {
	literal := &ArrayLiteral{
		Token:    *p.current,
		Elements: p.parseExpressionList(lexer.RightBracket),
	}

	return literal
}

func (p *Parser) parseHashLiteral() Expression {
	return nil
}

func (p *Parser) tryPeek(tokenToken lexer.TokenType) bool {
	if p.next.Type == tokenToken {
		p.nextToken()
		return true
	}

	p.errors = append(p.errors,
		NewParseError(
			fmt.Sprintf("expected %s, got %s", tokenToken, p.next.Type), p.l.GetSource(), p.next))

	return false
}

func (p *Parser) nextToken() {

	p.prev = p.current
	p.current = p.next

	next, err := p.l.Read()
	if err != nil {
		return
	}

	p.next = &next
}

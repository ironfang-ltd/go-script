package template

import (
	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
	"slices"
)

type Node interface {
}

type TextNode struct {
	Value string
}

type IfNode struct {
	Condition *parser.Program
	Body      []Node
	Else      []Node
}

type ForNode struct {
	Variable   string
	Collection string
	Body       []Node
}

type Parser struct {
	lexer   *Lexer
	current lexer.Token
	stack   []*lexer.TokenType
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (p *Parser) Parse() (*parser.Program, error) {

	err := p.consume()
	if err != nil {
		return nil, err
	}

	return p.ParseUntil([]lexer.TokenType{lexer.EndOfFile})
}

func (p *Parser) ParseUntil(until []lexer.TokenType) (*parser.Program, error) {

	program := parser.NewProgram()

	for {
		token, err := p.lexer.Read()
		if err != nil {
			return nil, err
		}

		if token.Type == lexer.EndOfFile {
			break
		}

		if slices.Contains(until, token.Type) {
			return program, nil
		}

		switch token.Type {
		case lexer.Text:
			// TODO: Add print call to program instead of TextNode
			//program.AddNode(&TextNode{Value: token.Value})

			err = p.consume()
			if err != nil {
				return nil, err
			}

		case lexer.Code:
			// Parse statement, if, foreach, etc.
			// parse the template into nodes and then walk the nodes to generate the statements?
		}
	}

	return program, nil
}

func (p *Parser) consume() error {
	current, err := p.lexer.Read()
	if err != nil {
		return err
	}

	p.current = current

	return nil
}

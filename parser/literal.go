package parser

import "github.com/ironfang-ltd/ironscript/lexer"

type IntegerLiteral struct {
	Token lexer.Token
	Value int
}

func (il *IntegerLiteral) Debug() string {
	return il.Token.Source
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) Debug() string {
	return sl.Token.Source
}

type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) Debug() string {
	return bl.Token.Source
}

type FunctionLiteral struct {
	Token      lexer.Token
	Identifier *Identifier
	Body       *BlockStatement
	Parameters []*Identifier
}

func (fl *FunctionLiteral) Debug() string {
	return fl.Token.Source + " " + fl.Identifier.Value + "(" + fl.Parameters[0].Value + ") " + fl.Body.Debug()
}

type ArrayLiteral struct {
	Token    lexer.Token
	Elements []Expression
}

func (al *ArrayLiteral) Debug() string {
	return al.Token.Source
}

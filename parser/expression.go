package parser

import "github.com/ironfang-ltd/ironscript/lexer"

type Expression interface {
	Debug() string
}

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) Debug() string {
	return i.Token.Source
}

type PropertyExpression struct {
	Token    lexer.Token
	Left     Expression
	Property Expression
}

func (pe *PropertyExpression) Debug() string {
	return pe.Left.Debug() + "." + pe.Property.Debug()
}

type CallExpression struct {
	Token    lexer.Token
	Function Expression
	Args     []Expression
}

func (ce *CallExpression) Debug() string {
	args := ""
	for _, arg := range ce.Args {
		args += arg.Debug() + ","
	}
	return ce.Function.Debug() + "(" + args + ")"
}

type InfixExpression struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (ie *InfixExpression) Debug() string {
	return ie.Left.Debug() + " " + ie.Token.Source + " " + ie.Right.Debug()
}

type IndexExpression struct {
	Token lexer.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Debug() string {
	return ie.Left.Debug() + "[" + ie.Index.Debug() + "]"
}

type IfExpression struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) Debug() string {
	alternative := ""
	if ie.Alternative != nil {
		alternative = " else " + ie.Alternative.Debug()
	}
	return "if " + ie.Condition.Debug() + " " + ie.Consequence.Debug() + alternative
}

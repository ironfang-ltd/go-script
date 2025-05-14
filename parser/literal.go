package parser

import (
	"strings"

	"github.com/ironfang-ltd/go-script/lexer"
)

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
	var sb strings.Builder
	sb.WriteString(fl.Token.Source)
	sb.WriteString(" ")
	sb.WriteString(fl.Identifier.Value)
	sb.WriteString("(")
	for i, p := range fl.Parameters {
		sb.WriteString(p.Value)
		if i < len(fl.Parameters)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") \n")
	sb.WriteString(fl.Body.Debug())
	return sb.String()
}

type ArrayLiteral struct {
	Token    lexer.Token
	Elements []Expression
}

func (al *ArrayLiteral) Debug() string {
	return al.Token.Source
}

type HashLiteral struct {
	Token lexer.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) Debug() string {
	var sb strings.Builder
	sb.WriteString(hl.Token.Source)
	sb.WriteString("{")
	for k, v := range hl.Pairs {
		sb.WriteString(k.Debug())
		sb.WriteString(": ")
		sb.WriteString(v.Debug())
	}
	sb.WriteString("}")
	return sb.String()
}

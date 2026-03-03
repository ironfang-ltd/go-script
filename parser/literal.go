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
	if fl.Identifier != nil {
		sb.WriteString(" ")
		sb.WriteString(fl.Identifier.Value)
	}
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

type NullLiteral struct {
	Token lexer.Token
}

func (nl *NullLiteral) Debug() string {
	return "null"
}

type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) Debug() string {
	return fl.Token.Source
}

type ArrayLiteral struct {
	Token    lexer.Token
	Elements []Expression
}

func (al *ArrayLiteral) Debug() string {
	return al.Token.Source
}

type HashPair struct {
	Key   Expression
	Value Expression
}

type HashLiteral struct {
	Token lexer.Token
	Pairs []HashPair
}

func (hl *HashLiteral) Debug() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, pair := range hl.Pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(pair.Key.Debug())
		sb.WriteString(": ")
		sb.WriteString(pair.Value.Debug())
	}
	sb.WriteString("}")
	return sb.String()
}

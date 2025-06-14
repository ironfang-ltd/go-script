package parser

import "github.com/ironfang-ltd/go-script/lexer"

type Statement interface {
	Debug() string
}

type ScriptStatement struct {
	Statements []Statement
}

func (ss *ScriptStatement) Debug() string {
	str := ""
	for _, s := range ss.Statements {
		str += s.Debug() + "\n"
	}
	return str
}

type PrintStatement struct {
	Value string
}

func (ps *PrintStatement) Debug() string {
	return "print " + ps.Value
}

type AssignmentExpression struct {
	Token lexer.Token
	Left  Expression
	Right Expression
}

func (ls *AssignmentExpression) Debug() string {
	return ls.Left.Debug() + " = " + ls.Right.Debug()
}

type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) Debug() string {
	return ls.Token.Source + " " + ls.Name.Value + " = " + ls.Value.Debug()
}

type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) Debug() string {
	return rs.Token.Source + " " + rs.Value.Debug()
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) Debug() string {
	return es.Expression.Debug()
}

type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) Debug() string {
	str := "{\n"
	for _, s := range bs.Statements {
		str += "    " + s.Debug() + "\n"
	}
	str += "}"
	return str
}

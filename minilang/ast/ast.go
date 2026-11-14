package ast

import "minilang/lexer"

type Node interface {
	GetToken() lexer.Token
}

type NodeBase struct {
	Token lexer.Token
}

func (n NodeBase) GetToken() lexer.Token { return n.Token }

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	NodeBase
	Name       *lexer.Token
	Statements []Statement
}

type VarDecl struct {
	NodeBase
	Name lexer.Token
	Type lexer.Token
}

type AssignStmt struct {
	NodeBase
	Name  lexer.Token
	Value Expression
}

type IfStmt struct {
	NodeBase
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

type WhileStmt struct {
	NodeBase
	Condition Expression
	Body      Statement
}

type PrintStmt struct {
	NodeBase
	Expression Expression
}

type ReadStmt struct {
	NodeBase
	Name lexer.Token
}

type BlockStmt struct {
	NodeBase
	Statements []Statement
}

type ExpressionStmt struct {
	NodeBase
	Expression Expression
}

type BinaryExpr struct {
	NodeBase
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

type UnaryExpr struct {
	NodeBase
	Operator lexer.Token
	Right    Expression
}

type GroupingExpr struct {
	NodeBase
	Expression Expression
}

type LiteralExpr struct {
	NodeBase
	Value interface{}
}

type IdentifierExpr struct {
	NodeBase
	Name lexer.Token
}

func (*VarDecl) statementNode()         {}
func (*AssignStmt) statementNode()      {}
func (*IfStmt) statementNode()          {}
func (*WhileStmt) statementNode()       {}
func (*PrintStmt) statementNode()       {}
func (*ReadStmt) statementNode()        {}
func (*BlockStmt) statementNode()       {}
func (*ExpressionStmt) statementNode()  {}
func (*BinaryExpr) expressionNode()     {}
func (*UnaryExpr) expressionNode()      {}
func (*GroupingExpr) expressionNode()   {}
func (*LiteralExpr) expressionNode()    {}
func (*IdentifierExpr) expressionNode() {}

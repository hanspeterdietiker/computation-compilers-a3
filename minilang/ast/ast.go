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
	GetType() Type
}

type Type string

const (
	TypeUnknown Type = "unknown"
	TypeInteger Type = "inteiro"
	TypeReal    Type = "real"
	TypeBoolean Type = "booleano"
)

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
	Left         Expression
	Operator     lexer.Token
	Right        Expression
	InferredType Type
}

type UnaryExpr struct {
	NodeBase
	Operator     lexer.Token
	Right        Expression
	InferredType Type
}

type GroupingExpr struct {
	NodeBase
	Expression   Expression
	InferredType Type
}

type LiteralExpr struct {
	NodeBase
	Value        interface{}
	InferredType Type
}

type IdentifierExpr struct {
	NodeBase
	Name         lexer.Token
	InferredType Type
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

func (e *BinaryExpr) GetType() Type   { return e.InferredType }
func (e *UnaryExpr) GetType() Type    { return e.InferredType }
func (e *GroupingExpr) GetType() Type { return e.InferredType }
func (e *LiteralExpr) GetType() Type  { return e.InferredType }
func (e *IdentifierExpr) GetType() Type {
	return e.InferredType
}

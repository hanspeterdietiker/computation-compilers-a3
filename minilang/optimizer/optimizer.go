// Package optimizer aplica dobra de constantes e eliminação de código morto na AST validada.
package optimizer

import (
	"minilang/ast"
	"minilang/lexer"
	"minilang/value"
)

// Report conta expressões dobradas e comandos se/enquanto eliminados.
type Report struct {
	FoldedExpressions int
	RemovedStatements int
}

type optimizer struct {
	report Report
}

// Optimize altera o programa no lugar; exige AST anotada pelo semântico.
func Optimize(program *ast.Program) Report {
	o := &optimizer{}
	if program != nil {
		program.Statements = o.statements(program.Statements)
	}
	return o.report
}

func (o *optimizer) statements(statements []ast.Statement) []ast.Statement {
	result := statements[:0]
	for _, statement := range statements {
		if optimized := o.statement(statement); optimized != nil {
			result = append(result, optimized)
		}
	}
	return result
}

// statement devolve nil quando o comando pode ser removido.
func (o *optimizer) statement(statement ast.Statement) ast.Statement {
	switch node := statement.(type) {
	case *ast.AssignStmt:
		node.Value = o.expression(node.Value)
	case *ast.PrintStmt:
		node.Expression = o.expression(node.Expression)
	case *ast.ExpressionStmt:
		node.Expression = o.expression(node.Expression)
	case *ast.BlockStmt:
		node.Statements = o.statements(node.Statements)
	case *ast.IfStmt:
		return o.ifStatement(node)
	case *ast.WhileStmt:
		node.Condition = o.expression(node.Condition)
		if isFalse(node.Condition) && !isBareDeclaration(node.Body) {
			o.report.RemovedStatements++
			return nil
		}
		node.Body = o.statement(node.Body)
		if node.Body == nil {
			node.Body = &ast.BlockStmt{NodeBase: ast.NodeBase{Token: node.Token}}
		}
	}
	return statement
}

func (o *optimizer) ifStatement(node *ast.IfStmt) ast.Statement {
	node.Condition = o.expression(node.Condition)
	if condition, constant := boolLiteral(node.Condition); constant {
		taken, discarded := node.ThenBranch, node.ElseBranch
		if !condition {
			taken, discarded = discarded, taken
		}
		// Declaração sem bloco pertence ao escopo externo: não remover.
		if !isBareDeclaration(discarded) {
			o.report.RemovedStatements++
			if taken == nil {
				return nil
			}
			return o.statement(taken)
		}
	}

	node.ThenBranch = o.statement(node.ThenBranch)
	if node.ThenBranch == nil {
		node.ThenBranch = &ast.BlockStmt{NodeBase: ast.NodeBase{Token: node.Token}}
	}
	if node.ElseBranch != nil {
		node.ElseBranch = o.statement(node.ElseBranch)
	}
	return node
}

func (o *optimizer) expression(expression ast.Expression) ast.Expression {
	switch node := expression.(type) {
	case *ast.GroupingExpr:
		node.Expression = o.expression(node.Expression)
		if literal, ok := node.Expression.(*ast.LiteralExpr); ok {
			return o.fold(node, literal.Value)
		}
	case *ast.UnaryExpr:
		node.Right = o.expression(node.Right)
		if operand, ok := node.Right.(*ast.LiteralExpr); ok {
			if result, err := value.Unary(node.Operator.Type, operand.Value); err == nil {
				return o.fold(node, result)
			}
		}
	case *ast.BinaryExpr:
		return o.binary(node)
	}
	return expression
}

func (o *optimizer) binary(node *ast.BinaryExpr) ast.Expression {
	node.Left = o.expression(node.Left)
	node.Right = o.expression(node.Right)

	// falso e X / verdadeiro ou X
	if left, constant := boolLiteral(node.Left); constant {
		if (node.Operator.Type == lexer.TokenAnd && !left) ||
			(node.Operator.Type == lexer.TokenOr && left) {
			return o.fold(node, left)
		}
	}

	left, leftIsLiteral := node.Left.(*ast.LiteralExpr)
	right, rightIsLiteral := node.Right.(*ast.LiteralExpr)
	if !leftIsLiteral || !rightIsLiteral {
		return node
	}
	result, err := value.Binary(node.Operator.Type, left.Value, right.Value)
	if err != nil { // divisão por zero fica para a execução
		return node
	}
	return o.fold(node, result)
}

// fold troca a expressão por um literal com o mesmo token e tipo.
func (o *optimizer) fold(original ast.Expression, result interface{}) ast.Expression {
	o.report.FoldedExpressions++
	typ := original.GetType()
	return &ast.LiteralExpr{
		NodeBase:     ast.NodeBase{Token: original.GetToken()},
		Value:        value.Coerce(typ, result),
		InferredType: typ,
	}
}

func boolLiteral(expression ast.Expression) (bool, bool) {
	literal, ok := expression.(*ast.LiteralExpr)
	if !ok {
		return false, false
	}
	b, ok := literal.Value.(bool)
	return b, ok
}

func isFalse(expression ast.Expression) bool {
	b, constant := boolLiteral(expression)
	return constant && !b
}

func isBareDeclaration(statement ast.Statement) bool {
	_, ok := statement.(*ast.VarDecl)
	return ok
}

package semantic

import (
	"fmt"

	"minilang/ast"
	"minilang/lexer"
)

type Error struct {
	Token   lexer.Token
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("[ERRO SEMÂNTICO] linha %d, coluna %d: %s", e.Token.Line, e.Token.Column, e.Message)
}

type Analyzer struct {
	symbols *SymbolTable
	errors  []error
}

func New() *Analyzer {
	return &Analyzer{symbols: NewSymbolTable()}
}

func (a *Analyzer) SymbolTable() *SymbolTable {
	return a.symbols
}

func (a *Analyzer) Analyze(program *ast.Program) []error {
	a.symbols = NewSymbolTable()
	a.errors = nil
	if program == nil {
		a.report(lexer.Token{}, "programa ausente")
		return append([]error(nil), a.errors...)
	}
	for _, statement := range program.Statements {
		a.analyzeStatement(statement)
	}
	return append([]error(nil), a.errors...)
}

func (a *Analyzer) analyzeStatement(statement ast.Statement) {
	switch node := statement.(type) {
	case *ast.VarDecl:
		if _, declared := a.symbols.Declare(node.Name.Lexeme, typeFromToken(node.Type), node.Name); !declared {
			a.report(node.Name, fmt.Sprintf("variável '%s' já declarada no escopo '%s'.", node.Name.Lexeme, a.symbols.CurrentScope().Name))
		}
	case *ast.AssignStmt:
		symbol, declared := a.symbols.Resolve(node.Name.Lexeme)
		if !declared {
			a.report(node.Name, fmt.Sprintf("variável '%s' não declarada.", node.Name.Lexeme))
		}
		valueType := a.analyzeExpression(node.Value)
		if declared && valueType != ast.TypeUnknown && !assignable(symbol.Type, valueType) {
			a.report(node.Value.GetToken(), fmt.Sprintf(
				"não é possível atribuir um valor do tipo %s a uma variável do tipo %s.",
				valueType, symbol.Type,
			))
		}
	case *ast.IfStmt:
		a.requireBoolean(node.Condition, "condição do 'se'")
		a.analyzeStatement(node.ThenBranch)
		if node.ElseBranch != nil {
			a.analyzeStatement(node.ElseBranch)
		}
	case *ast.WhileStmt:
		a.requireBoolean(node.Condition, "condição do 'enquanto'")
		a.analyzeStatement(node.Body)
	case *ast.PrintStmt:
		a.analyzeExpression(node.Expression)
	case *ast.ReadStmt:
		if _, declared := a.symbols.Resolve(node.Name.Lexeme); !declared {
			a.report(node.Name, fmt.Sprintf("variável '%s' não declarada em 'leia'.", node.Name.Lexeme))
		}
	case *ast.BlockStmt:
		a.symbols.PushScope(node.Token)
		for _, child := range node.Statements {
			a.analyzeStatement(child)
		}
		a.symbols.PopScope()
	case *ast.ExpressionStmt:
		a.analyzeExpression(node.Expression)
	}
}

func (a *Analyzer) requireBoolean(expression ast.Expression, description string) {
	typ := a.analyzeExpression(expression)
	if typ != ast.TypeUnknown && typ != ast.TypeBoolean {
		a.report(expression.GetToken(), fmt.Sprintf("%s deve ser booleana, mas recebeu %s.", description, typ))
	}
}

func (a *Analyzer) analyzeExpression(expression ast.Expression) ast.Type {
	if expression == nil {
		return ast.TypeUnknown
	}
	switch node := expression.(type) {
	case *ast.LiteralExpr:
		switch node.Value.(type) {
		case int64:
			node.InferredType = ast.TypeInteger
		case float64:
			node.InferredType = ast.TypeReal
		case bool:
			node.InferredType = ast.TypeBoolean
		default:
			node.InferredType = ast.TypeUnknown
		}
	case *ast.IdentifierExpr:
		symbol, declared := a.symbols.Resolve(node.Name.Lexeme)
		if !declared {
			a.report(node.Name, fmt.Sprintf("variável '%s' não declarada.", node.Name.Lexeme))
			node.InferredType = ast.TypeUnknown
		} else {
			node.InferredType = symbol.Type
		}
	case *ast.GroupingExpr:
		node.InferredType = a.analyzeExpression(node.Expression)
	case *ast.UnaryExpr:
		rightType := a.analyzeExpression(node.Right)
		switch node.Operator.Type {
		case lexer.TokenNot:
			node.InferredType = ast.TypeBoolean
			if rightType != ast.TypeUnknown && rightType != ast.TypeBoolean {
				a.report(node.Operator, fmt.Sprintf("operador '%s' exige operando booleano, mas recebeu %s.", node.Operator.Lexeme, rightType))
				node.InferredType = ast.TypeUnknown
			}
		case lexer.TokenPlus, lexer.TokenMinus:
			node.InferredType = rightType
			if rightType != ast.TypeUnknown && !isNumeric(rightType) {
				a.report(node.Operator, fmt.Sprintf("operador '%s' exige operando numérico, mas recebeu %s.", node.Operator.Lexeme, rightType))
				node.InferredType = ast.TypeUnknown
			}
		}
	case *ast.BinaryExpr:
		leftType := a.analyzeExpression(node.Left)
		rightType := a.analyzeExpression(node.Right)
		node.InferredType = a.binaryType(node.Operator, leftType, rightType)
	default:
		return ast.TypeUnknown
	}
	return expression.GetType()
}

func (a *Analyzer) binaryType(operator lexer.Token, left, right ast.Type) ast.Type {
	if left == ast.TypeUnknown || right == ast.TypeUnknown {
		return ast.TypeUnknown
	}
	switch operator.Type {
	case lexer.TokenPlus, lexer.TokenMinus, lexer.TokenMultiply, lexer.TokenDivide:
		if isNumeric(left) && isNumeric(right) {
			if left == ast.TypeReal || right == ast.TypeReal {
				return ast.TypeReal
			}
			return ast.TypeInteger
		}
		a.report(operator, fmt.Sprintf("operador aritmético '%s' exige operandos numéricos, mas recebeu %s e %s.", operator.Lexeme, left, right))
	case lexer.TokenPercent:
		if left == ast.TypeInteger && right == ast.TypeInteger {
			return ast.TypeInteger
		}
		a.report(operator, fmt.Sprintf("operador '%%' exige operandos inteiros, mas recebeu %s e %s.", left, right))
	case lexer.TokenAnd, lexer.TokenOr:
		if left == ast.TypeBoolean && right == ast.TypeBoolean {
			return ast.TypeBoolean
		}
		a.report(operator, fmt.Sprintf("operador lógico '%s' exige operandos booleanos, mas recebeu %s e %s.", operator.Lexeme, left, right))
	case lexer.TokenLessThan, lexer.TokenLessEqual, lexer.TokenGreaterThan, lexer.TokenGreaterEqual:
		if isNumeric(left) && isNumeric(right) {
			return ast.TypeBoolean
		}
		a.report(operator, fmt.Sprintf("operador relacional '%s' exige operandos numéricos, mas recebeu %s e %s.", operator.Lexeme, left, right))
	case lexer.TokenEqual, lexer.TokenNotEqual:
		if (isNumeric(left) && isNumeric(right)) || left == right {
			return ast.TypeBoolean
		}
		a.report(operator, fmt.Sprintf("operador de igualdade '%s' exige operandos compatíveis, mas recebeu %s e %s.", operator.Lexeme, left, right))
	}
	return ast.TypeUnknown
}

func (a *Analyzer) report(token lexer.Token, message string) {
	a.errors = append(a.errors, Error{Token: token, Message: message})
}

func typeFromToken(token lexer.Token) ast.Type {
	switch token.Type {
	case lexer.TokenTypeInt:
		return ast.TypeInteger
	case lexer.TokenTypeFloat:
		return ast.TypeReal
	case lexer.TokenTypeBool:
		return ast.TypeBoolean
	default:
		return ast.TypeUnknown
	}
}

func assignable(target, value ast.Type) bool {
	return target == value || target == ast.TypeReal && value == ast.TypeInteger
}

func isNumeric(typ ast.Type) bool {
	return typ == ast.TypeInteger || typ == ast.TypeReal
}

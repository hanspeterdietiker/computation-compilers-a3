// Package interpreter executa a AST anotada pelo analisador semântico.
package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"minilang/ast"
	"minilang/lexer"
	"minilang/value"
)

// RuntimeError é um erro de execução posicionado no token de origem.
type RuntimeError struct {
	Token   lexer.Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[ERRO DE EXECUÇÃO] linha %d, coluna %d: %s", e.Token.Line, e.Token.Column, e.Message)
}

// Interpreter lê de input (leia) e escreve em output (escreva).
type Interpreter struct {
	input  *bufio.Scanner
	output io.Writer
	env    *environment
}

// New cria um interpretador; cada leia consome uma palavra da entrada.
func New(input io.Reader, output io.Writer) *Interpreter {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)
	return &Interpreter{input: scanner, output: output}
}

// Run executa o programa e para no primeiro erro de execução.
func (i *Interpreter) Run(program *ast.Program) error {
	i.env = newEnvironment(nil)
	if program == nil {
		return nil
	}
	return i.executeAll(program.Statements)
}

func (i *Interpreter) executeAll(statements []ast.Statement) error {
	for _, statement := range statements {
		if err := i.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(statement ast.Statement) error {
	switch node := statement.(type) {
	case *ast.VarDecl:
		typ := ast.TypeFromToken(node.Type)
		i.env.declare(node.Name.Lexeme, typ, value.Zero(typ))
	case *ast.AssignStmt:
		result, err := i.evaluate(node.Value)
		if err != nil {
			return err
		}
		return i.assign(node.Name, result)
	case *ast.IfStmt:
		i.declareBareBranch(node.ThenBranch)
		i.declareBareBranch(node.ElseBranch)
		condition, err := i.evaluate(node.Condition)
		if err != nil {
			return err
		}
		if condition.(bool) {
			return i.execute(node.ThenBranch)
		}
		if node.ElseBranch != nil {
			return i.execute(node.ElseBranch)
		}
	case *ast.WhileStmt:
		i.declareBareBranch(node.Body)
		for {
			condition, err := i.evaluate(node.Condition)
			if err != nil {
				return err
			}
			if !condition.(bool) {
				return nil
			}
			if err := i.execute(node.Body); err != nil {
				return err
			}
		}
	case *ast.PrintStmt:
		result, err := i.evaluate(node.Expression)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(i.output, value.Format(result)); err != nil {
			return RuntimeError{Token: node.Token, Message: fmt.Sprintf("falha ao escrever a saída: %v", err)}
		}
	case *ast.ReadStmt:
		return i.read(node)
	case *ast.BlockStmt:
		previous := i.env
		i.env = newEnvironment(previous)
		defer func() { i.env = previous }()
		return i.executeAll(node.Statements)
	case *ast.ExpressionStmt:
		_, err := i.evaluate(node.Expression)
		return err
	}
	return nil
}

// declareBareBranch declara um var sem bloco no escopo externo, como no semântico.
func (i *Interpreter) declareBareBranch(branch ast.Statement) {
	if declaration, ok := branch.(*ast.VarDecl); ok {
		_ = i.execute(declaration)
	}
}

func (i *Interpreter) read(node *ast.ReadStmt) error {
	variable, ok := i.env.resolve(node.Name.Lexeme)
	if !ok {
		return RuntimeError{Token: node.Name, Message: fmt.Sprintf("variável '%s' não declarada.", node.Name.Lexeme)}
	}
	if !i.input.Scan() {
		message := fmt.Sprintf("entrada insuficiente para 'leia(%s)'.", node.Name.Lexeme)
		if err := i.input.Err(); err != nil {
			message = fmt.Sprintf("falha ao ler a entrada: %v", err)
		}
		return RuntimeError{Token: node.Token, Message: message}
	}
	parsed, err := value.Parse(variable.typ, i.input.Text())
	if err != nil {
		return RuntimeError{Token: node.Token, Message: err.Error() + "."}
	}
	variable.value = parsed
	return nil
}

func (i *Interpreter) assign(name lexer.Token, result interface{}) error {
	variable, ok := i.env.resolve(name.Lexeme)
	if !ok {
		return RuntimeError{Token: name, Message: fmt.Sprintf("variável '%s' não declarada.", name.Lexeme)}
	}
	variable.value = value.Coerce(variable.typ, result)
	return nil
}

func (i *Interpreter) evaluate(expression ast.Expression) (interface{}, error) {
	switch node := expression.(type) {
	case *ast.LiteralExpr:
		return node.Value, nil
	case *ast.IdentifierExpr:
		variable, ok := i.env.resolve(node.Name.Lexeme)
		if !ok {
			return nil, RuntimeError{Token: node.Name, Message: fmt.Sprintf("variável '%s' não declarada.", node.Name.Lexeme)}
		}
		return variable.value, nil
	case *ast.GroupingExpr:
		return i.evaluate(node.Expression)
	case *ast.UnaryExpr:
		operand, err := i.evaluate(node.Right)
		if err != nil {
			return nil, err
		}
		result, err := value.Unary(node.Operator.Type, operand)
		return i.wrap(node.Operator, result, err)
	case *ast.BinaryExpr:
		return i.binary(node)
	}
	return nil, RuntimeError{Token: expression.GetToken(), Message: fmt.Sprintf("expressão não suportada: %T", expression)}
}

func (i *Interpreter) binary(node *ast.BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	// curto-circuito
	switch node.Operator.Type {
	case lexer.TokenAnd:
		if !left.(bool) {
			return false, nil
		}
	case lexer.TokenOr:
		if left.(bool) {
			return true, nil
		}
	}
	right, err := i.evaluate(node.Right)
	if err != nil {
		return nil, err
	}
	result, err := value.Binary(node.Operator.Type, left, right)
	return i.wrap(node.Operator, result, err)
}

// wrap posiciona erros das operações no operador.
func (i *Interpreter) wrap(operator lexer.Token, result interface{}, err error) (interface{}, error) {
	if err == nil {
		return result, nil
	}
	message := err.Error()
	if errors.Is(err, value.ErrDivisionByZero) {
		message = "divisão por zero."
	}
	return nil, RuntimeError{Token: operator, Message: message}
}

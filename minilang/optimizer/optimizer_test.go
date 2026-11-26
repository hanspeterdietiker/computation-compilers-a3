package optimizer

import (
	"testing"

	"minilang/ast"
	"minilang/lexer"
	"minilang/parser"
	"minilang/semantic"
)

func TestOptimizeFoldsConstantExpressionsKeepingTypes(t *testing.T) {
	program := analyzed(t, `
		var total: real;
		var n: inteiro;
		total = 10 + 2.5 * 2;
		n = -(3 + 4) * 2;
		escreva(não falso e (1 < 2));
	`)
	report := Optimize(program)

	assertLiteral(t, program.Statements[2].(*ast.AssignStmt).Value, 15.0, ast.TypeReal)
	assertLiteral(t, program.Statements[3].(*ast.AssignStmt).Value, int64(-14), ast.TypeInteger)
	assertLiteral(t, program.Statements[4].(*ast.PrintStmt).Expression, true, ast.TypeBoolean)
	if report.FoldedExpressions == 0 {
		t.Error("relatório não contabilizou expressões dobradas")
	}
}

func TestOptimizeFoldsOnlyConstantSubexpressions(t *testing.T) {
	program := analyzed(t, `
		var x: inteiro;
		x = x + 2 * 3;
	`)
	Optimize(program)
	sum, ok := program.Statements[1].(*ast.AssignStmt).Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expressão com variável não deveria ser dobrada: %T", program.Statements[1].(*ast.AssignStmt).Value)
	}
	assertLiteral(t, sum.Right, int64(6), ast.TypeInteger)
}

func TestOptimizeKeepsDivisionByZeroForRuntime(t *testing.T) {
	program := analyzed(t, `escreva(1 / 0);`)
	Optimize(program)
	if _, ok := program.Statements[0].(*ast.PrintStmt).Expression.(*ast.BinaryExpr); !ok {
		t.Fatal("divisão por zero não deve ser dobrada")
	}
}

func TestOptimizeShortCircuitsLeftConstant(t *testing.T) {
	program := analyzed(t, `
		var a: booleano;
		escreva(falso e a);
		escreva(verdadeiro ou a);
	`)
	Optimize(program)
	assertLiteral(t, program.Statements[1].(*ast.PrintStmt).Expression, false, ast.TypeBoolean)
	assertLiteral(t, program.Statements[2].(*ast.PrintStmt).Expression, true, ast.TypeBoolean)
}

func TestOptimizeEliminatesDeadBranches(t *testing.T) {
	program := analyzed(t, `
		var x: inteiro;
		se (1 > 2) escreva(1); senão escreva(2);
		se (falso) escreva(3);
		enquanto (falso) { x = x + 1; }
		se (2 > 1) { escreva(4); }
	`)
	report := Optimize(program)

	if len(program.Statements) != 3 {
		t.Fatalf("comandos após otimização = %d; esperados 3", len(program.Statements))
	}
	assertLiteral(t, program.Statements[1].(*ast.PrintStmt).Expression, int64(2), ast.TypeInteger)
	if _, ok := program.Statements[2].(*ast.BlockStmt); !ok {
		t.Fatalf("ramo verdadeiro deveria substituir o 'se': %T", program.Statements[2])
	}
	if report.RemovedStatements != 4 {
		t.Errorf("comandos removidos = %d; esperados 4", report.RemovedStatements)
	}
}

func TestOptimizeKeepsBareDeclarationBranch(t *testing.T) {
	program := analyzed(t, `
		se (falso) var y: inteiro;
		y = 1;
	`)
	Optimize(program)
	if _, ok := program.Statements[0].(*ast.IfStmt); !ok {
		t.Fatalf("'se' com declaração direta não deve ser removido: %T", program.Statements[0])
	}
}

func analyzed(t *testing.T, source string) *ast.Program {
	t.Helper()
	program, errors := parser.New(lexer.New(source).ScanTokens()).Parse()
	if len(errors) != 0 {
		t.Fatalf("erros sintáticos inesperados: %v", errors)
	}
	if errors := semantic.New().Analyze(program); len(errors) != 0 {
		t.Fatalf("erros semânticos inesperados: %v", errors)
	}
	return program
}

func assertLiteral(t *testing.T, expression ast.Expression, want interface{}, typ ast.Type) {
	t.Helper()
	literal, ok := expression.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("expressão = %T; esperado literal %v", expression, want)
	}
	if literal.Value != want || literal.GetType() != typ {
		t.Fatalf("literal = %#v (%s); esperado %#v (%s)", literal.Value, literal.GetType(), want, typ)
	}
}

package semantic

import (
	"strings"
	"testing"

	"minilang/ast"
	"minilang/lexer"
	"minilang/parser"
)

func TestAnalyzeAcceptsRealPromotionAndAnnotatesExpressions(t *testing.T) {
	program := parseProgram(t, `
		var preco: real;
		var quantidade: inteiro;
		preco = 19.90;
		quantidade = 3;
		preco = 10;
		preco = preco * quantidade;
	`)
	analyzer := New()
	if errors := analyzer.Analyze(program); len(errors) != 0 {
		t.Fatalf("erros semânticos inesperados: %v", errors)
	}
	promotedAssignment := program.Statements[4].(*ast.AssignStmt)
	if got := promotedAssignment.Value.GetType(); got != ast.TypeInteger {
		t.Errorf("tipo da expressão promovida = %s; esperado %s", got, ast.TypeInteger)
	}
	assignment := program.Statements[5].(*ast.AssignStmt)
	expression := assignment.Value.(*ast.BinaryExpr)
	if got := expression.GetType(); got != ast.TypeReal {
		t.Errorf("tipo inferido = %s; esperado %s", got, ast.TypeReal)
	}
	if got := expression.Right.GetType(); got != ast.TypeInteger {
		t.Errorf("tipo do operando inteiro = %s; esperado %s", got, ast.TypeInteger)
	}
	if got := len(analyzer.SymbolTable().Symbols()); got != 2 {
		t.Errorf("símbolos declarados = %d; esperados 2", got)
	}
}

func TestAnalyzeReportsSemanticErrorsAndPositions(t *testing.T) {
	program := parseProgram(t, `var count: inteiro;
count = 1.5;
missing = count + verdadeiro;
se (count) escreva(count);
leia(unknown);`)
	errors := New().Analyze(program)
	if len(errors) != 5 {
		t.Fatalf("erros semânticos = %d; esperados 5: %v", len(errors), errors)
	}
	for _, semanticError := range errors {
		if !strings.HasPrefix(semanticError.Error(), "[ERRO SEMÂNTICO] linha ") {
			t.Errorf("diagnóstico sem fase/localização: %q", semanticError)
		}
	}
	if got := errors[0].Error(); !strings.Contains(got, "linha 2, coluna 9") || !strings.Contains(got, "real") {
		t.Errorf("erro de atribuição incorreta = %q", got)
	}
	if got := errors[1].Error(); !strings.Contains(got, "variável 'missing' não declarada") {
		t.Errorf("erro de variável não declarada = %q", got)
	}
	if got := errors[2].Error(); !strings.Contains(got, "operador aritmético") {
		t.Errorf("erro de operador = %q", got)
	}
	if got := errors[3].Error(); !strings.Contains(got, "deve ser booleana") {
		t.Errorf("erro de condição = %q", got)
	}
	if got := errors[4].Error(); !strings.Contains(got, "variável 'unknown' não declarada em 'leia'") {
		t.Errorf("erro de leia = %q", got)
	}
}

func TestSymbolTableSupportsNestedScopesAndShadowing(t *testing.T) {
	program := parseProgram(t, `
		var item: inteiro;
		{
			var item: real;
			item = 2.5;
		}
		item = 1;
		escreva(item);
	`)
	analyzer := New()
	if errors := analyzer.Analyze(program); len(errors) != 0 {
		t.Fatalf("erros semânticos inesperados: %v", errors)
	}
	symbols := analyzer.SymbolTable().Symbols()
	if len(symbols) != 2 || symbols[0].Scope == symbols[1].Scope {
		t.Fatalf("símbolos/escopos inesperados: %+v", symbols)
	}
	if symbols[0].Type != ast.TypeInteger || symbols[1].Type != ast.TypeReal {
		t.Errorf("tipos dos símbolos = %s, %s", symbols[0].Type, symbols[1].Type)
	}
	if symbols[0].DeclarationToken.Line == 0 || symbols[0].DeclarationToken.Column == 0 {
		t.Errorf("posição de declaração não registrada: %+v", symbols[0].DeclarationToken)
	}
	print := program.Statements[3].(*ast.PrintStmt)
	if got := print.Expression.GetType(); got != ast.TypeInteger {
		t.Errorf("tipo resolvido após sair do bloco = %s; esperado %s", got, ast.TypeInteger)
	}
}

func TestAnalyzeDetectsRedeclarationAndInvalidOperators(t *testing.T) {
	program := parseProgram(t, `
		var flag: booleano;
		var flag: booleano;
		flag = verdadeiro + 1;
		flag = não 1;
	`)
	errors := New().Analyze(program)
	if len(errors) != 3 {
		t.Fatalf("erros semânticos = %d; esperados 3: %v", len(errors), errors)
	}
	if !strings.Contains(errors[0].Error(), "já declarada") ||
		!strings.Contains(errors[1].Error(), "operador aritmético") ||
		!strings.Contains(errors[2].Error(), "operador 'não'") {
		t.Errorf("diagnósticos inesperados: %v", errors)
	}
}

func parseProgram(t *testing.T, source string) *ast.Program {
	t.Helper()
	program, errors := parser.New(lexer.New(source).ScanTokens()).Parse()
	if len(errors) != 0 {
		t.Fatalf("erros sintáticos inesperados: %v", errors)
	}
	return program
}

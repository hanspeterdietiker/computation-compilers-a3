package interpreter

import (
	"strings"
	"testing"

	"minilang/ast"
	"minilang/lexer"
	"minilang/optimizer"
	"minilang/parser"
	"minilang/semantic"
)

func TestRunExecutesArithmeticWithRealCoercion(t *testing.T) {
	output := mustRun(t, `
		var a: inteiro;
		var b: real;
		var resultado: real;
		a = 10;
		b = 2.5;
		resultado = a + b;
		escreva(resultado);
		b = 10;
		escreva(b);
		escreva(a / 4);
		escreva(a % 4);
		escreva(-a * 2);
	`, "")
	assertOutput(t, output, "12.5", "10.0", "2", "2", "-20")
}

func TestRunDeclarationsStartWithZeroValues(t *testing.T) {
	output := mustRun(t, `
		var i: inteiro; var r: real; var b: booleano;
		escreva(i); escreva(r); escreva(b);
	`, "")
	assertOutput(t, output, "0", "0.0", "falso")
}

func TestRunConditionalsAndLoops(t *testing.T) {
	output := mustRun(t, `
		var i: inteiro;
		var soma: inteiro;
		i = 1;
		enquanto (i <= 5) {
			soma = soma + i;
			se (i % 2 == 0) escreva(i); senão { }
			i = i + 1;
		}
		se (soma == 15 e não (soma < 0)) escreva(verdadeiro); senão escreva(falso);
		escreva(soma);
	`, "")
	assertOutput(t, output, "2", "4", "verdadeiro", "15")
}

func TestRunBlocksCreateScopesWithShadowing(t *testing.T) {
	output := mustRun(t, `
		var x: inteiro;
		x = 1;
		{
			var x: real;
			x = 2.5;
			escreva(x);
		}
		escreva(x);
		{ x = 7; }
		escreva(x);
	`, "")
	assertOutput(t, output, "2.5", "1", "7")
}

func TestRunFreshScopeOnEachLoopIteration(t *testing.T) {
	output := mustRun(t, `
		var i: inteiro;
		enquanto (i < 2) {
			var local: inteiro;
			escreva(local);
			local = 9;
			i = i + 1;
		}
	`, "")
	assertOutput(t, output, "0", "0")
}

func TestRunReadsTypedInput(t *testing.T) {
	output := mustRun(t, `
		var n: inteiro; var r: real; var b: booleano;
		leia(n); leia(r); leia(b);
		escreva(n * 2); escreva(r); escreva(não b);
	`, "21\n3\nfalso\n")
	assertOutput(t, output, "42", "3.0", "verdadeiro")
}

func TestRunShortCircuitsLogicalOperators(t *testing.T) {
	output := mustRun(t, `
		var zero: inteiro;
		se (zero != 0 e 10 / zero > 1) escreva(1); senão escreva(2);
		se (zero == 0 ou 10 / zero > 1) escreva(3);
	`, "")
	assertOutput(t, output, "2", "3")
}

func TestRunReportsRuntimeErrors(t *testing.T) {
	cases := []struct {
		name, source, input, message string
	}{
		{"divisão por zero", "var z: inteiro;\nescreva(1 / z);", "", "linha 2, coluna 11: divisão por zero"},
		{"módulo por zero", "var z: inteiro;\nescreva(1 % z);", "", "divisão por zero"},
		{"entrada ausente", "var n: inteiro;\nleia(n);", "", "linha 2, coluna 1: entrada insuficiente"},
		{"entrada inválida", "var n: inteiro;\nleia(n);", "abc", "valor 'abc' inválido para inteiro"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output strings.Builder
			err := New(strings.NewReader(tc.input), &output).Run(analyzed(t, tc.source))
			if err == nil {
				t.Fatal("erro de execução esperado")
			}
			if !strings.HasPrefix(err.Error(), "[ERRO DE EXECUÇÃO] ") || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("erro = %q; esperado conter %q", err, tc.message)
			}
		})
	}
}

func TestRunStopsAtFirstRuntimeErrorKeepingPreviousOutput(t *testing.T) {
	var output strings.Builder
	err := New(strings.NewReader(""), &output).Run(analyzed(t, `
		var z: inteiro;
		escreva(1);
		escreva(1 / z);
		escreva(2);
	`))
	if err == nil {
		t.Fatal("erro de execução esperado")
	}
	assertOutput(t, output.String(), "1")
}

func TestRunOptimizedProgramProducesSameOutput(t *testing.T) {
	source := `
		var r: real;
		var i: inteiro;
		r = 10 + 2.5 * 2;
		se (1 < 2 e verdadeiro) escreva(r); senão escreva(0);
		enquanto (falso) escreva(99);
		enquanto (i < 3) { i = i + 2 * 1; }
		escreva(i);
	`
	plain := mustRun(t, source, "")
	program := analyzed(t, source)
	optimizer.Optimize(program)
	var optimized strings.Builder
	if err := New(strings.NewReader(""), &optimized).Run(program); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if plain != optimized.String() {
		t.Fatalf("saída otimizada = %q; sem otimização = %q", optimized.String(), plain)
	}
	assertOutput(t, plain, "15.0", "4")
}

func mustRun(t *testing.T, source, input string) string {
	t.Helper()
	var output strings.Builder
	if err := New(strings.NewReader(input), &output).Run(analyzed(t, source)); err != nil {
		t.Fatalf("erro de execução inesperado: %v", err)
	}
	return output.String()
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

func assertOutput(t *testing.T, got string, lines ...string) {
	t.Helper()
	want := strings.Join(lines, "\n") + "\n"
	if got != want {
		t.Fatalf("saída = %q; esperada %q", got, want)
	}
}

func TestRunDeclaresBareBranchDeclarationRegardlessOfCondition(t *testing.T) {
	output := mustRun(t, `
		se (falso) var y: inteiro;
		y = 1;
		escreva(y);
	`, "")
	assertOutput(t, output, "1")
}

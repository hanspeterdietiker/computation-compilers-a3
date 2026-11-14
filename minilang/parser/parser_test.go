package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"minilang/ast"
	"minilang/lexer"
)

func TestParserBuildsPrecedenceAndRealLiteralAST(t *testing.T) {
	program := parseValid(t, "var total: float; total = 10 + 2.5 * 3;")
	if len(program.Statements) != 2 {
		t.Fatalf("quantidade de comandos = %d; esperada 2", len(program.Statements))
	}

	assignment, ok := program.Statements[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("comando de atribuição tem tipo %T", program.Statements[1])
	}

	addition, ok := assignment.Value.(*ast.BinaryExpr)
	if !ok || addition.Operator.Type != lexer.TokenPlus {
		t.Fatalf("raiz da expressão = %#v; esperada adição", assignment.Value)
	}
	multiplication, ok := addition.Right.(*ast.BinaryExpr)
	if !ok || multiplication.Operator.Type != lexer.TokenMultiply {
		t.Fatalf("operando direito da adição = %#v; esperada multiplicação", addition.Right)
	}
	floatLiteral, ok := multiplication.Left.(*ast.LiteralExpr)
	if !ok {
		t.Fatalf("operando esquerdo da multiplicação tem tipo %T", multiplication.Left)
	}
	if _, ok := floatLiteral.Value.(float64); !ok {
		t.Fatalf("literal real armazenado como %T", floatLiteral.Value)
	}
}

func TestParserUsesLeftAssociativity(t *testing.T) {
	program := parseValid(t, "var resultado: int; resultado = 20 - 5 - 2;")
	assignment := program.Statements[1].(*ast.AssignStmt)
	root := assignment.Value.(*ast.BinaryExpr)
	left, ok := root.Left.(*ast.BinaryExpr)
	if !ok || left.Operator.Type != lexer.TokenMinus || root.Operator.Type != lexer.TokenMinus {
		t.Fatalf("subtrações não foram agrupadas à esquerda: %#v", assignment.Value)
	}
}

func TestParserBindsElseToNearestIf(t *testing.T) {
	program := parseValid(t, "if (true) if (false) print(1); else print(2);")
	outer := program.Statements[0].(*ast.IfStmt)
	inner, ok := outer.ThenBranch.(*ast.IfStmt)
	if !ok {
		t.Fatalf("ramo then externo tem tipo %T", outer.ThenBranch)
	}
	if outer.ElseBranch != nil || inner.ElseBranch == nil {
		t.Fatal("else não foi associado ao if mais próximo")
	}
}

func TestParserRecoversAndReportsMultipleErrors(t *testing.T) {
	program, errors := parse("var a int; print(1); var b: int; print(2;")
	if len(errors) != 2 {
		t.Fatalf("erros sintáticos = %d; esperados 2: %v", len(errors), errors)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("comandos válidos preservados = %d; esperados 2", len(program.Statements))
	}
}

func TestParserAcceptsPortugueseKeywords(t *testing.T) {
	program := parseValid(t, "programa exemplo { var ativo: booleano; ativo = verdadeiro; se (não falso) escreva(ativo); }")
	if program.Name == nil || program.Name.Lexeme != "exemplo" {
		t.Fatalf("nome do programa = %v; esperado exemplo", program.Name)
	}
	if len(program.Statements) != 3 {
		t.Fatalf("quantidade de comandos = %d; esperada 3", len(program.Statements))
	}
}

func TestParserBuildsControlFlowAndIOAST(t *testing.T) {
	program := parseValid(t, "var x: int; read(x); while (x < 3) { x = x + 1; print(x); }")
	if len(program.Statements) != 3 {
		t.Fatalf("quantidade de comandos = %d; esperada 3", len(program.Statements))
	}
	if _, ok := program.Statements[1].(*ast.ReadStmt); !ok {
		t.Fatalf("comando read tem tipo %T", program.Statements[1])
	}

	loop, ok := program.Statements[2].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("comando while tem tipo %T", program.Statements[2])
	}
	body, ok := loop.Body.(*ast.BlockStmt)
	if !ok || len(body.Statements) != 2 {
		t.Fatalf("corpo do while = %#v; esperado bloco com 2 comandos", loop.Body)
	}
	if _, ok := body.Statements[0].(*ast.AssignStmt); !ok {
		t.Fatalf("primeiro comando do bloco tem tipo %T", body.Statements[0])
	}
	if _, ok := body.Statements[1].(*ast.PrintStmt); !ok {
		t.Fatalf("segundo comando do bloco tem tipo %T", body.Statements[1])
	}
}

func TestParserReportsUnexpectedTokenPositionAndMessage(t *testing.T) {
	_, errors := parse("var item int;")
	if len(errors) != 1 {
		t.Fatalf("erros sintáticos = %d; esperado 1: %v", len(errors), errors)
	}
	parseError, ok := errors[0].(ParseError)
	if !ok {
		t.Fatalf("tipo de erro = %T; esperado ParseError", errors[0])
	}
	if parseError.Token.Line != 1 || parseError.Token.Column != 10 {
		t.Errorf("posição do erro = %d:%d; esperada 1:10", parseError.Token.Line, parseError.Token.Column)
	}
	want := "[ERRO SINTÁTICO] linha 1, coluna 10: esperado ':' após o nome da variável"
	if parseError.Error() != want {
		t.Errorf("mensagem = %q; esperada %q", parseError.Error(), want)
	}
}

func TestValidProgramCorpus(t *testing.T) {
	paths, err := filepath.Glob("testdata/valid/*.min")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("nenhum programa válido encontrado em testdata/valid")
	}

	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			program, errors := parse(string(source))
			if len(errors) != 0 {
				t.Fatalf("programa válido rejeitado: %v", errors)
			}
			if len(program.Statements) == 0 {
				t.Fatal("AST do programa válido não contém comandos")
			}
		})
	}
}

func TestInvalidProgramCorpus(t *testing.T) {
	paths, err := filepath.Glob("testdata/invalid/*.min")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("nenhum programa inválido encontrado em testdata/invalid")
	}

	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			wantBytes, err := os.ReadFile(strings.TrimSuffix(path, ".min") + ".want")
			if err != nil {
				t.Fatal(err)
			}
			_, errors := parse(string(source))
			got := make([]string, len(errors))
			for i, parseError := range errors {
				got[i] = parseError.Error()
			}
			normalizedWant := strings.ReplaceAll(string(wantBytes), "\r\n", "\n")
			want := strings.Split(strings.TrimSpace(normalizedWant), "\n")
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Fatalf("diagnósticos:\n%s\nesperados:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
		})
	}
}

func parseValid(t *testing.T, source string) *ast.Program {
	t.Helper()
	program, errors := parse(source)
	if len(errors) != 0 {
		t.Fatalf("erros inesperados: %v", errors)
	}
	return program
}

func parse(source string) (*ast.Program, []error) {
	return New(lexer.New(source).ScanTokens()).Parse()
}

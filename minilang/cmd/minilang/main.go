package main

// para testar rode 'go run ./cmd/minilang run examples/programa.min'
// os testes automatizados são executados com 'go test ./...'

import (
	"fmt"
	"io"
	"os"

	"minilang/ast"
	"minilang/interpreter"
	"minilang/lexer"
	"minilang/optimizer"
	"minilang/parser"
	"minilang/semantic"
)

const usage = "Uso: minilang <arquivo> | minilang parse <arquivo> | minilang check <arquivo> | minilang opt <arquivo> | minilang run <arquivo>"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, input io.Reader, output, errorOutput io.Writer) int {
	mode := "tokens"
	if len(args) > 0 {
		switch args[0] {
		case "parse", "check", "opt", "run":
			mode, args = args[0], args[1:]
		}
	}
	if len(args) != 1 {
		if mode == "tokens" {
			fmt.Fprintln(errorOutput, usage)
		} else {
			fmt.Fprintf(errorOutput, "Uso: minilang %s <arquivo>\n", mode)
		}
		return 2
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintln(errorOutput, "Erro ao ler arquivo:", err)
		return 2
	}

	scanner := lexer.New(string(data))
	tokens := scanner.ScanTokens()
	if mode == "tokens" {
		printTokens(output, tokens)
		return 0
	}
	if len(scanner.Errors()) > 0 {
		fmt.Fprintf(errorOutput, "Análise léxica encerrada: %d erro(s) encontrado(s).\n", len(scanner.Errors()))
		return 1
	}

	program, errors := parser.New(tokens).Parse()
	if reportErrors(errorOutput, errors, "sintática") {
		return 1
	}
	if mode == "parse" {
		fmt.Fprintf(output, "Análise sintática concluída sem erros. AST criada com %d comando(s).\n", len(program.Statements))
		return 0
	}

	analyzer := semantic.New()
	if reportErrors(errorOutput, analyzer.Analyze(program), "semântica") {
		return 1
	}
	switch mode {
	case "check":
		printSymbols(output, analyzer.SymbolTable())
		return 0
	case "opt":
		report := optimizer.Optimize(program)
		fmt.Fprintf(output, "Otimização concluída: %d expressão(ões) constante(s) dobrada(s), %d comando(s) eliminado(s).\n",
			report.FoldedExpressions, report.RemovedStatements)
		return 0
	}
	return execute(program, input, output, errorOutput)
}

func execute(program *ast.Program, input io.Reader, output, errorOutput io.Writer) int {
	optimizer.Optimize(program)
	if err := interpreter.New(input, output).Run(program); err != nil {
		fmt.Fprintln(errorOutput, err)
		return 1
	}
	return 0
}

func reportErrors(errorOutput io.Writer, errors []error, phase string) bool {
	for _, err := range errors {
		fmt.Fprintln(errorOutput, err)
	}
	if len(errors) > 0 {
		fmt.Fprintf(errorOutput, "Análise %s encerrada: %d erro(s) encontrado(s).\n", phase, len(errors))
	}
	return len(errors) > 0
}

func printSymbols(output io.Writer, table *semantic.SymbolTable) {
	symbols := table.Symbols()
	fmt.Fprintf(output, "Análise semântica concluída sem erros. AST anotada com %d símbolo(s).\n", len(symbols))
	fmt.Fprintln(output, "Identificador | Tipo | Escopo | Declaração")
	for _, symbol := range symbols {
		fmt.Fprintf(output, "%s | %s | %s | %d:%d\n",
			symbol.Name, symbol.Type, symbol.Scope,
			symbol.DeclarationToken.Line, symbol.DeclarationToken.Column)
	}
}

func printTokens(output io.Writer, tokens []lexer.Token) {
	for _, token := range tokens {
		fmt.Fprintf(output, "%-20s %-15q %d:%d\n", token.Type, token.Lexeme, token.Line, token.Column)
	}
}

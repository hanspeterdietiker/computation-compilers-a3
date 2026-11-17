package main

// existe o package exemples
// para testar rode 'go run ./cmd/minilang examples/exemple-identifer.min'
// para rodar deve ter o go instalado e configurado no path do sistema operacional
// além do mais há testes automatizados que podem ser executados com o comando 'go test'

import (
	"fmt"
	"io"
	"os"

	"minilang/lexer"
	"minilang/parser"
	"minilang/semantic"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, output, errorOutput io.Writer) int {
	parseMode := len(args) > 0 && args[0] == "parse"
	checkMode := len(args) > 0 && args[0] == "check"
	if parseMode || checkMode {
		args = args[1:]
	}
	if len(args) != 1 {
		if parseMode {
			fmt.Fprintln(errorOutput, "Uso: minilang parse <arquivo>")
		} else if checkMode {
			fmt.Fprintln(errorOutput, "Uso: minilang check <arquivo>")
		} else {
			fmt.Fprintln(errorOutput, "Uso: minilang <arquivo> | minilang parse <arquivo> | minilang check <arquivo>")
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
	if (parseMode || checkMode) && len(scanner.Errors()) > 0 {
		fmt.Fprintf(errorOutput, "Análise léxica encerrada: %d erro(s) encontrado(s).\n", len(scanner.Errors()))
		return 1
	}
	if parseMode || checkMode {
		program, errors := parser.New(tokens).Parse()
		for _, parseError := range errors {
			fmt.Fprintln(errorOutput, parseError)
		}
		if len(errors) > 0 {
			fmt.Fprintf(errorOutput, "Análise sintática encerrada: %d erro(s) encontrado(s).\n", len(errors))
			return 1
		}
		if checkMode {
			analyzer := semantic.New()
			semanticErrors := analyzer.Analyze(program)
			for _, semanticError := range semanticErrors {
				fmt.Fprintln(errorOutput, semanticError)
			}
			if len(semanticErrors) > 0 {
				fmt.Fprintf(errorOutput, "Análise semântica encerrada: %d erro(s) encontrado(s).\n", len(semanticErrors))
				return 1
			}
			fmt.Fprintf(output, "Análise semântica concluída sem erros. AST anotada com %d símbolo(s).\n", len(analyzer.SymbolTable().Symbols()))
			fmt.Fprintln(output, "Identificador | Tipo | Escopo | Declaração")
			for _, symbol := range analyzer.SymbolTable().Symbols() {
				fmt.Fprintf(output, "%s | %s | %s | %d:%d\n",
					symbol.Name, symbol.Type, symbol.Scope,
					symbol.DeclarationToken.Line, symbol.DeclarationToken.Column)
			}
			return 0
		}
		fmt.Fprintf(output, "Análise sintática concluída sem erros. AST criada com %d comando(s).\n", len(program.Statements))
		return 0
	}

	for _, token := range tokens {
		fmt.Fprintf(
			output,
			"%-20s %-15q %d:%d\n",
			token.Type,
			token.Lexeme,
			token.Line,
			token.Column,
		)
	}
	return 0
}

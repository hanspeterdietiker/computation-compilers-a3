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
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, output, errorOutput io.Writer) int {
	parseMode := len(args) > 0 && args[0] == "parse"
	if parseMode {
		args = args[1:]
	}
	if len(args) != 1 {
		if parseMode {
			fmt.Fprintln(errorOutput, "Uso: minilang parse <arquivo>")
		} else {
			fmt.Fprintln(errorOutput, "Uso: minilang <arquivo> | minilang parse <arquivo>")
		}
		return 2
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintln(errorOutput, "Erro ao ler arquivo:", err)
		return 2
	}

	tokens := lexer.New(string(data)).ScanTokens()
	if parseMode {
		program, errors := parser.New(tokens).Parse()
		for _, parseError := range errors {
			fmt.Fprintln(errorOutput, parseError)
		}
		if len(errors) > 0 {
			fmt.Fprintf(errorOutput, "Análise sintática encerrada: %d erro(s) encontrado(s).\n", len(errors))
			return 1
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

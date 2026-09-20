package main

import (
	"fmt"
	"os"

	"minilang/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: minilang <arquivo>")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Erro ao ler arquivo:", err)
		return
	}

	l := lexer.New(string(data))

	tokens := l.ScanTokens()

	for _, token := range tokens {
		fmt.Printf(
			"%-20s %-15q %d:%d\n",
			token.Type,
			token.Lexeme,
			token.Line,
			token.Column,
		)
	}
}

package main

// existe o package exemples
// para testar rode 'go run ./cmd/minilang examples/exemple-identifer.min'
// para rodar deve ter o go instalado e configurado no path do sistema operacional
// além do mais há testes automatizados que podem ser executados com o comando 'go test'

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

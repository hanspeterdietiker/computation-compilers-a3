package lexer

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestIntegerEmitsIntLiteral(t *testing.T) {
	tokens := New("20").ScanTokens()

	assertTokens(t, tokens, []expectedToken{
		{TokenIntLiteral, "20", 1, 1},
		{TokenEOF, "", 1, 3},
	})
}

func TestHashCommentIsDiscarded(t *testing.T) {
	tokens := New("var idade # comentário\n42").ScanTokens()

	assertTokens(t, tokens, []expectedToken{
		{TokenVar, "var", 1, 1},
		{TokenIdentifier, "idade", 1, 5},
		{TokenIntLiteral, "42", 2, 1},
		{TokenEOF, "", 2, 3},
	})
}

func TestSlashRemainsDivisionOperator(t *testing.T) {
	tokens := New("//").ScanTokens()

	assertTokens(t, tokens, []expectedToken{
		{TokenDivide, "/", 1, 1},
		{TokenDivide, "/", 1, 2},
		{TokenEOF, "", 1, 3},
	})
}

func TestLoneExclamationReportsLexicalError(t *testing.T) {
	output := captureStdout(t, func() {
		tokens := New("!").ScanTokens()
		assertTokens(t, tokens, []expectedToken{
			{TokenEOF, "", 1, 2},
		})
	})

	want := "[ERRO LÉXICO] linha 1, coluna 1: caractere '!' não reconhecido"
	if !strings.Contains(output, want) {
		t.Fatalf("saída de erro = %q; esperado que contenha %q", output, want)
	}
}

type expectedToken struct {
	tokenType TokenType
	lexeme    string
	line      int
	column    int
}

func assertTokens(t *testing.T, got []Token, want []expectedToken) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("quantidade de tokens = %d; esperada %d", len(got), len(want))
	}

	for i, expected := range want {
		actual := got[i]
		if actual.Type != expected.tokenType ||
			actual.Lexeme != expected.lexeme ||
			actual.Line != expected.line ||
			actual.Column != expected.column {
			t.Errorf(
				"token[%d] = {%s %q %d:%d}; esperado {%s %q %d:%d}",
				i,
				actual.Type,
				actual.Lexeme,
				actual.Line,
				actual.Column,
				expected.tokenType,
				expected.lexeme,
				expected.line,
				expected.column,
			)
		}
	}
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("não foi possível criar pipe para stdout: %v", err)
	}

	os.Stdout = writer
	defer func() {
		os.Stdout = original
	}()

	run()

	if err := writer.Close(); err != nil {
		t.Fatalf("não foi possível fechar stdout temporário: %v", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("não foi possível ler stdout temporário: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("não foi possível fechar a leitura de stdout: %v", err)
	}

	return string(output)
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Casos: .min, saída em .want, entrada opcional em .in, erro esperado em .err.
func TestRunPrograms(t *testing.T) {
	sources, err := filepath.Glob(filepath.Join("testdata", "run", "*.min"))
	if err != nil || len(sources) == 0 {
		t.Fatalf("nenhum programa de teste encontrado: %v", err)
	}
	for _, source := range sources {
		base := strings.TrimSuffix(source, ".min")
		t.Run(filepath.Base(base), func(t *testing.T) {
			want := readFile(t, base+".want")
			wantErr, expectFailure := readOptional(base + ".err")
			input, _ := readOptional(base + ".in")

			var output, errorOutput strings.Builder
			code := run([]string{"run", source}, strings.NewReader(input), &output, &errorOutput)

			if got := output.String(); got != want {
				t.Errorf("saída = %q; esperada %q", got, want)
			}
			if expectFailure {
				if code != 1 || !strings.Contains(errorOutput.String(), strings.TrimSpace(wantErr)) {
					t.Errorf("código %d, erros %q; esperado código 1 contendo %q", code, errorOutput.String(), wantErr)
				}
			} else if code != 0 || errorOutput.Len() != 0 {
				t.Errorf("código %d, erros %q; esperada execução sem erros", code, errorOutput.String())
			}
		})
	}
}

func TestOptModeReportsOptimizations(t *testing.T) {
	var output, errorOutput strings.Builder
	code := run([]string{"opt", filepath.Join("testdata", "run", "otimizacao.min")}, strings.NewReader(""), &output, &errorOutput)
	if code != 0 {
		t.Fatalf("código %d: %s", code, errorOutput.String())
	}
	want := "Otimização concluída: 7 expressão(ões) constante(s) dobrada(s), 3 comando(s) eliminado(s).\n"
	if output.String() != want {
		t.Errorf("saída = %q; esperada %q", output.String(), want)
	}
}

func TestRunExamples(t *testing.T) {
	for _, name := range []string{"programa.min", "real.min"} {
		var output, errorOutput strings.Builder
		path := filepath.Join("..", "..", "examples", name)
		if code := run([]string{"run", path}, strings.NewReader("4\n"), &output, &errorOutput); code != 0 {
			t.Errorf("%s: código %d: %s", name, code, errorOutput.String())
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("falha ao ler %s: %v", path, err)
	}
	return normalize(string(data))
}

func readOptional(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return normalize(string(data)), true
}

func normalize(text string) string { return strings.ReplaceAll(text, "\r\n", "\n") }

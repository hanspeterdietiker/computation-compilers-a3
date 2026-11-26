# Guia de uso

## Pré-requisitos

* Go **1.27.1** ou superior (versão declarada em `minilang/go.mod`).
* Todos os comandos abaixo rodam **dentro de `minilang/`**:

```bash
cd minilang
```

Para gerar um executável em vez de usar `go run`:

```bash
go build -o minilang ./cmd/minilang     # no Windows: -o minilang.exe
./minilang run examples/real.min
```

## Comandos

```text
minilang <arquivo>          mostra os tokens (só a análise léxica)
minilang parse <arquivo>    léxica + sintática
minilang check <arquivo>    léxica + sintática + semântica, mostra a tabela de símbolos
minilang opt <arquivo>      tudo acima + otimização, mostra o que foi otimizado
minilang run <arquivo>      pipeline completo: executa o programa
```

Cada modo executa todas as fases anteriores. Se uma fase encontra erros, as
seguintes não rodam.

### `run`: executar um programa

```bash
$ go run ./cmd/minilang run examples/real.min
59.7
10.0
12.5
```

`leia` consome valores da entrada padrão, separados por espaço ou quebra de
linha:

```bash
$ echo 6 | go run ./cmd/minilang run examples/programa.min
2
4
6
21
7.0
13
12.5
```

### `check`: ver a tabela de símbolos

```bash
$ go run ./cmd/minilang check examples/real.min
Análise semântica concluída sem erros. AST anotada com 3 símbolo(s).
Identificador | Tipo | Escopo | Declaração
preco | real | global | 1:5
quantidade | inteiro | global | 2:5
total | real | global | 3:5
```

### `opt`: ver o efeito da otimização

```bash
$ go run ./cmd/minilang opt examples/programa.min
Otimização concluída: 5 expressão(ões) constante(s) dobrada(s), 0 comando(s) eliminado(s).
```

### `parse`: validar a sintaxe

```bash
$ go run ./cmd/minilang parse examples/parser/valido.min
Análise sintática concluída sem erros. AST criada com 3 comando(s).
```

### Sem modo: listar os tokens

Cada linha mostra o tipo do token, o lexema e a posição `linha:coluna`. O
último token é sempre `EOF`.

```bash
$ go run ./cmd/minilang arquivo.min      # contendo: var x: inteiro;
VAR                  "var"           1:1
IDENTIFIER           "x"             1:5
COLON                ":"             1:6
TYPE_INT             "inteiro"       1:8
SEMICOLON            ";"             1:15
EOF                  ""              1:16
```

## Mensagens de erro

Toda mensagem informa a fase, a linha, a coluna e o problema:

```text
[ERRO LÉXICO] linha 7, coluna 9: caractere '@' não reconhecido
[ERRO SINTÁTICO] linha 1, coluna 14: esperado ':' após o nome da variável
[ERRO SEMÂNTICO] linha 2, coluna 1: variável 'y' não declarada.
[ERRO SEMÂNTICO] linha 2, coluna 5: não é possível atribuir um valor do tipo real a uma variável do tipo inteiro.
[ERRO SEMÂNTICO] linha 3, coluna 5: condição do 'se' deve ser booleana, mas recebeu inteiro.
[ERRO DE EXECUÇÃO] linha 5, coluna 11: divisão por zero.
[ERRO DE EXECUÇÃO] linha 2, coluna 1: valor 'dez' inválido para inteiro.
```

* As fases léxica, sintática e semântica relatam **todos** os erros que
  encontram e terminam com um resumo, por exemplo
  `Análise sintática encerrada: 2 erro(s) encontrado(s).`
* A execução para no **primeiro** erro, mas mantém o que já foi impresso.

## Códigos de saída

| Código | Significado |
| --- | --- |
| `0` | Sucesso |
| `1` | Erro léxico, sintático, semântico ou de execução |
| `2` | Uso incorreto da CLI ou arquivo não encontrado |

O modo de tokens sempre termina com `0`; os erros léxicos aparecem junto com a
lista de tokens.

## Exemplos incluídos

| Arquivo | O que demonstra |
| --- | --- |
| `examples/programa.min` | Programa completo: `leia`, laço, condicional, média em `real`, otimização |
| `examples/real.min` | Extensão C: operações mistas e atribuição `inteiro → real` |
| `examples/exemple.min` | Tokens de operadores e delimitadores |
| `examples/exemple-identifer.min` | Identificadores, números e um erro léxico proposital (`@`) |
| `examples/parser/valido.min` | Programa sintaticamente válido |
| `examples/parser/erro-declaracao.min` | Erro sintático em declaração |
| `examples/parser/multiplos-erros.min` | Dois erros sintáticos relatados na mesma execução |

## Testes

```bash
go test ./...            # todos os pacotes
go test -v ./parser      # lista cada caso da bateria do parser
```

A organização dos testes está descrita em
[`arquitetura.md`](../tecnico/arquitetura.md#testes).

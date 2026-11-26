# Back-end (M4)

| | |
| --- | --- |
| **Entrada** | `*ast.Program` anotado pelo [analisador semântico](semantico.md) |
| **Saída** | Saída padrão do programa, ou um erro de execução |
| **Estratégia** | Otimização sobre a AST seguida de interpretação direta da AST |
| **Código** | `minilang/value/`, `minilang/optimizer/`, `minilang/interpreter/` |
| **API** | `report := optimizer.Optimize(program)`; `err := interpreter.New(stdin, stdout).Run(program)` |

O back-end recebe a AST validada e anotada pelo analisador semântico, aplica a
otimização e executa o programa por **interpretação direta da AST**. Não há
geração de código de máquina nem de bytecode: o interpretador percorre os nós
da árvore e executa cada um.

```text
AST anotada → Otimizador → AST otimizada → Interpretador → saída
```

## Pacotes

| Pacote        | Responsabilidade                                                        |
| ------------- | ----------------------------------------------------------------------- |
| `value`       | Valores (`int64`, `float64`, `bool`), operadores, coerção e formatação. |
| `optimizer`   | Dobra de constantes e eliminação de código morto.                        |
| `interpreter` | Execução dos comandos, ambientes de escopo, `leia` e `escreva`.          |

O otimizador e o interpretador usam as mesmas operações de `value`, então uma
expressão dobrada em tempo de compilação produz exatamente o valor que seria
calculado na execução.

## Interpretador

* Cada declaração cria a variável com valor inicial `0`, `0.0` ou `falso`.
* Cada bloco cria um ambiente filho, espelhando os escopos do M3; o corpo de um
  `enquanto` recebe um ambiente novo a cada iteração.
* Atribuição a uma variável `real` converte valores `inteiro` (`10` → `10.0`).
* Em operações mistas, o operando `inteiro` é promovido para `real`.
* `/` entre inteiros é divisão inteira (trunca em direção a zero); `%` exige
  inteiros.
* `e` e `ou` são avaliados em curto-circuito.
* `leia` consome uma palavra da entrada padrão e a converte para o tipo da
  variável; `real` aceita a forma inteira e `booleano` aceita
  `verdadeiro`/`falso` (ou `true`/`false`).
* `escreva` imprime um valor por linha. Reais sempre mostram a parte decimal
  (`10.0`) e usam até 15 dígitos significativos (`19.90 * 3` → `59.7`).
* Uma declaração usada como ramo de `se`/`enquanto` sem bloco é declarada no
  escopo externo independentemente da condição, como no analisador semântico.

### Erros de execução

A execução para no primeiro erro, mantendo a saída já produzida:

```text
[ERRO DE EXECUÇÃO] linha 5, coluna 11: divisão por zero.
[ERRO DE EXECUÇÃO] linha 2, coluna 1: entrada insuficiente para 'leia(n)'.
[ERRO DE EXECUÇÃO] linha 2, coluna 1: valor 'dez' inválido para inteiro.
```

## Otimização

**Dobra de constantes.** Subexpressões formadas apenas por literais são
calculadas antes da execução e trocadas por um literal com o mesmo tipo
inferido:

```text
x = 2 * (3 + 4) - 1;   →   x = 13;
r = 10 + 2.5;          →   r = 12.5;
falso e a              →   falso
verdadeiro ou a        →   verdadeiro
```

Divisões por zero não são dobradas, para que o erro seja reportado na execução
com a posição original.

**Eliminação de código morto.** Depois da dobra:

```text
se (1 > 2) A; senão B;   →   B;
se (verdadeiro) A;       →   A;
se (falso) A;            →   (removido)
enquanto (falso) A;      →   (removido)
```

O modo `opt` da CLI mostra quantas expressões foram dobradas e quantos comandos
foram eliminados.

## Execução

Dentro de `minilang/`:

```bash
go run ./cmd/minilang run examples/programa.min
go run ./cmd/minilang opt examples/programa.min
go test ./value ./optimizer ./interpreter ./cmd/minilang
```

Os testes de integração ficam em `cmd/minilang/testdata/run/`: cada `.min` tem
a saída esperada em `.want`, entrada opcional em `.in` e, quando deve falhar, o
diagnóstico esperado em `.err`.

Como criar novos casos de integração está descrito em
[`arquitetura.md`](arquitetura.md#testes).

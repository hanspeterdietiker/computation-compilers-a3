# Análise semântica (M3)

| | |
| --- | --- |
| **Entrada** | `*ast.Program` sem erros sintáticos |
| **Saída** | A mesma AST com `InferredType` preenchido em cada expressão, a tabela de símbolos e `[]error` |
| **Modelo** | Tabela de símbolos com escopos aninhados e sistema de tipos |
| **Código** | `minilang/semantic/analyzer.go`, `minilang/semantic/symbol_table.go` |
| **API** | `a := semantic.New(); errs := a.Analyze(program); tabela := a.SymbolTable()` |

## Por que esta fase existe

A gramática livre de contexto só reconhece a **forma** do programa.
`x = verdadeiro + 1;` é sintaticamente perfeito. Regras como "a variável precisa
estar declarada", "os tipos precisam ser compatíveis" e "a condição precisa ser
booleana" dependem do contexto (o que foi declarado antes e com qual tipo) e
não podem ser expressas numa GLC. O analisador semântico verifica essas regras
percorrendo a AST.

## Tabela de símbolos e escopos

Cada símbolo guarda nome, tipo, escopo e o token da declaração (linha e
coluna).

* O escopo `global` contém as declarações do programa.
* Cada bloco `{ ... }` cria um escopo filho, removido ao terminar o bloco.
* A busca de um nome começa no escopo atual e sobe pelos ancestrais. Um bloco
  pode **sombrear** uma variável externa.
* Declarar o mesmo nome duas vezes no **mesmo** escopo é erro; a declaração
  repetida não substitui a original.
* Uma declaração só é visível **a partir do ponto** em que aparece.

## Regras de tipos

| Construção | Regra | Tipo resultante |
| --- | --- | --- |
| Literal | `10`, `2.5`, `verdadeiro` | `inteiro`, `real`, `booleano` |
| `+ - * /` | Operandos numéricos | `real` se algum for `real`; senão `inteiro` |
| `%` | Dois `inteiro` | `inteiro` |
| `< <= > >=` | Operandos numéricos | `booleano` |
| `== !=` | Dois numéricos ou dois `booleano` | `booleano` |
| `e ou não` | Operandos `booleano` | `booleano` |
| `-` / `+` unário | Operando numérico | Tipo do operando |
| Condição de `se` / `enquanto` | Precisa ser `booleano` | — |
| Atribuição | Tipos iguais, ou `inteiro` para `real` | — |
| `leia(x)` | `x` precisa estar declarado | — |

## Extensão C: promoção `inteiro → real`

É nesta fase que a extensão escolhida pela equipe é verificada:

* Em operações aritméticas e relacionais mistas, o operando `inteiro` é tratado
  como `real` e o resultado é `real`.
* A atribuição `real ← inteiro` é aceita; o back-end faz a conversão do valor.
* A atribuição `inteiro ← real` é rejeitada: não há conversão implícita que
  perca a parte decimal.

```text
var q: inteiro;
q = 5.75;
```

```text
[ERRO SEMÂNTICO] linha 2, coluna 5: não é possível atribuir um valor do tipo real a uma variável do tipo inteiro.
```

## Erros e AST anotada

A análise percorre todos os comandos e acumula os erros com linha e coluna:

```text
[ERRO SEMÂNTICO] linha 2, coluna 1: variável 'y' não declarada.
[ERRO SEMÂNTICO] linha 2, coluna 16: operador aritmético '+' exige operandos numéricos, mas recebeu booleano e inteiro.
[ERRO SEMÂNTICO] linha 3, coluna 5: condição do 'se' deve ser booleana, mas recebeu inteiro.
Análise semântica encerrada: 3 erro(s) encontrado(s).
```

Cada expressão recebe seu tipo em `InferredType`, consultado por `GetType()`.
Uma expressão com erro fica com o tipo `unknown` (`ast.TypeUnknown`), e as
regras que recebem `unknown` não geram novos erros. Isso evita que um único
problema produza uma cascata de mensagens.

O back-end usa os tipos anotados: o otimizador cria literais com o tipo
correto e o interpretador sabe quando converter `inteiro` para `real`.

## Execução e testes

```bash
cd minilang
go run ./cmd/minilang check examples/real.min   # imprime a tabela de símbolos
go test -v ./semantic
```

`analyzer_test.go` cobre a promoção para `real` e a anotação das expressões,
posições dos erros, escopos aninhados com sombreamento, redeclaração e
operadores inválidos.

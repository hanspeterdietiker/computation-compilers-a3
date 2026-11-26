# Análise sintática e AST (M2)

| | |
| --- | --- |
| **Entrada** | `[]lexer.Token` produzido pelo lexer |
| **Saída** | `*ast.Program` (AST completa ou parcial) e `[]error` |
| **Modelo formal** | Gramática livre de contexto ([`gramatica.ebnf`](gramatica.ebnf)) e autômato de pilha |
| **Código** | `minilang/parser/parser.go`, `minilang/ast/ast.go` |
| **API** | `program, errs := parser.New(tokens).Parse()` |

## Gramática

A sintaxe é definida pela GLC em [`gramatica.ebnf`](gramatica.ebnf). Os pontos
principais:

* As regras de expressão formam uma escada de precedência, da menor para a
  maior: `ou` → `e` → igualdade → comparação → `+ -` → `* / %` → unários →
  primários.
* Cada nível é `operando { operador operando }`, o que dá associatividade à
  esquerda: `a - b - c` é `(a - b) - c`.
* O **dangling else** é resolvido ligando o `senão` ao `se` mais interno que
  ainda não tem alternativa.

## Autômato de pilha

Toda GLC é reconhecida por um autômato de pilha não determinístico (AP). O AP da
MiniLang é obtido pela construção padrão a partir da gramática.

Seja `G = (N, T, P, S)` a GLC de `gramatica.ebnf`, com as regras de precedência
expandidas em não-terminais. O AP é a tupla

```text
M = (Q, T, Γ, δ, q₀, #, {qf})
Q = {q₀, q, qf}
Γ = N ∪ T ∪ {#}
```

`#` marca o fundo da pilha e o topo fica à esquerda. A entrada é a sequência de
tokens do lexer, terminada por `EOF`.

| # | Transição | Significado |
| --- | --- | --- |
| 1 | `δ(q₀, ε, #) ∋ (q, S#)` | Empilha o símbolo inicial |
| 2 | `δ(q, ε, A) ∋ (q, X₁X₂…Xₙ)` para cada produção `A → X₁X₂…Xₙ` | **Expande**: troca o não-terminal do topo pelo lado direito de uma produção |
| 3 | `δ(q, a, a) ∋ (q, ε)` para cada terminal `a` | **Casa**: consome o token se ele for igual ao terminal do topo |
| 4 | `δ(q, ε, #) ∋ (qf, ε)` | Depois de consumir `EOF`, aceita |

A palavra é aceita quando existe uma sequência de escolhas que consome todos os
tokens e chega a `qf`. O não determinismo está na transição 2, quando um
não-terminal tem mais de uma produção (por exemplo, as alternativas de
`statement` e de `primary`). Na transição 2, `X₁` fica no topo; uma
implementação que empilha um símbolo por vez deve empilhar o lado direito em
ordem inversa.

## Do autômato para o parser

O parser é **descendente recursivo**: a pilha do AP vira a pilha de chamadas
de funções do Go, e cada não-terminal vira uma função.

| Não-terminal | Função |
| --- | --- |
| `program` | `Parse()` |
| `statement` | `parseStatement()` |
| `variable_declaration` | `parseVarDecl()` |
| `conditional` | `parseIf()` |
| `loop` | `parseWhile()` |
| `output` / `input` | `parsePrint()` / `parseRead()` |
| `block` | `parseBlock()` |
| `logical_or` … `primary` | `parseOr()`, `parseAnd()`, `parseEquality()`, `parseComparison()`, `parseTerm()`, `parseFactor()`, `parseUnary()`, `parsePrimary()` |

O não determinismo do AP some porque o parser olha o próximo token
(*lookahead* de 1) para escolher a produção: `var` leva a uma declaração, `se`
a um condicional, um identificador a uma atribuição, e assim por diante. A
transição 3 corresponde a `consume(tipo, mensagem)`, que consome o token
esperado ou registra um erro.

## AST

Cada produção reconhecida gera um nó. Todos os nós guardam o token de origem
(`NodeBase.Token`), com linha e coluna, para que as fases seguintes relatem
erros na posição certa.

| Comandos (`ast.Statement`) | Expressões (`ast.Expression`) |
| --- | --- |
| `VarDecl`, `AssignStmt`, `IfStmt`, `WhileStmt`, `PrintStmt`, `ReadStmt`, `BlockStmt` | `BinaryExpr`, `UnaryExpr`, `GroupingExpr`, `LiteralExpr`, `IdentifierExpr` |

As expressões têm o campo `InferredType`, vazio após o parser e preenchido pelo
[analisador semântico](semantico.md).

Exemplo: `resultado = 10 + 5 * 2;`

```text
AssignStmt (resultado)
└── BinaryExpr (+)
    ├── LiteralExpr 10
    └── BinaryExpr (*)
        ├── LiteralExpr 5
        └── LiteralExpr 2
```

## Recuperação de erros

Quando uma produção falha, o parser registra o erro e entra em modo de
sincronização (`synchronize()`): descarta tokens até passar por um `;`, chegar
a um `}` ou encontrar o início de um comando (`var`, `se`, `enquanto`,
`escreva`, `leia`, `{` ou identificador). Assim ele consegue relatar vários
erros numa mesma execução:

```text
[ERRO SINTÁTICO] linha 1, coluna 14: esperado ':' após o nome da variável
[ERRO SINTÁTICO] linha 2, coluna 8: esperado ')' após expressão
Análise sintática encerrada: 2 erro(s) encontrado(s).
```

Internamente, a falha usa `panic(ParseError)`, que `parseStatementSafely()`
captura com `recover`.

## Testes

```bash
cd minilang
go test -v ./parser
```

* `parser_test.go`: precedência, associatividade, dangling else, palavras em
  português, recuperação de múltiplos erros e posição das mensagens.
* `parser/testdata/valid/`: programas que precisam ser aceitos.
* `parser/testdata/invalid/`: programas rejeitados, com os diagnósticos
  esperados no `.want` de mesmo nome.

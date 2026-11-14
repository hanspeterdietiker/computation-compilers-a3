# Automato de Pilha da MiniLang

A gramática da MiniLang em [`grammar.ebnf`](grammar.ebnf) e livre de contexto.
Seu reconhecimento pode ser descrito pelo autômato de pilha não determinístico
(AP/NPDA) obtido pela construção padrão de uma GLC.

## Definicao formal

Seja `G = (N, T, P, S)` a GLC documentada em `grammar.ebnf`, expandindo as
regras de precedência em não-terminais. O AP é a tupla
`M = (Q, T, Γ, δ, q₀, #, {qf})`, onde `Q = {q₀, q, qf}` e
`Γ = N ∪ T ∪ {#}`. `#` é o marcador de fundo e o topo da pilha fica à
esquerda. A entrada é a sequência de tokens do lexer, incluindo `EOF` como
sentinela terminal.

As transições são:

1. Inicialização: `δ(q₀, ε, #)` contém `(q, S#)`, colocando o símbolo inicial
   acima do marcador de fundo.
2. Para cada produção `A → X₁ X₂ ... Xₙ` em `P`,
   `δ(q, ε, A)` contém `(q, X₁ X₂ ... Xₙ)`. O autômato escolhe uma produção
   quando encontra um não-terminal no topo da pilha.
3. Para cada terminal `a` (palavra reservada, identificador, literal,
   operador ou delimitador), `δ(q, a, a)` contém `(q, ε)`. O token é consumido
   somente quando corresponde ao topo terminal.
4. Após consumir `EOF`, `δ(q, ε, #)` contém `(qf, ε)`. A palavra é aceita
   quando a entrada está vazia e o autômato chega a `qf`.

Na regra de expansão, a cadeia de substituição é empilhada com `X₁` no topo;
implementações cuja operação empilha um símbolo por vez devem inseri-la na
ordem inversa. A não determinização aparece nas escolhas de produções e nas
alternativas de `statement` e `primary`. Uma palavra é aceita se existe uma
sequência de escolhas que consome todos os tokens e esvazia a pilha.

## Relação com o parser

O parser implementado usa a mesma GLC de forma descendente recursiva. Os
procedimentos `parseOr`, `parseAnd`, `parseEquality`, `parseComparison`,
`parseTerm`, `parseFactor`, `parseUnary` e `parsePrimary` simulam a expansão
dos não-terminais de expressão. O lookahead escolhe a produção e elimina a
necessidade de explorar alternativas inválidas do AP. Ao reconhecer cada
produção, o parser constrói o nó AST correspondente. Em caso de erro, a
recuperação avança até um ponto de sincronização (`;`, `}` ou início de
comando), permitindo reportar mais de um erro na mesma análise.

Assim, o AP formaliza o poder de reconhecimento da gramática, enquanto a
implementação determinística produz também a árvore sintática abstrata.
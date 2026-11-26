# Arquitetura

O compilador é um pipeline de cinco fases. Cada fase é um pacote Go
independente que recebe a saída da fase anterior. Não há dependências
externas: só a biblioteca padrão.

```text
 código-fonte (.min)
        │
        ▼
 ┌──────────────┐  AFD                         lexico.md
 │    lexer     │  string → []lexer.Token
 └──────┬───────┘
        ▼
 ┌──────────────┐  GLC + autômato de pilha     gramatica.ebnf, sintatico.md
 │    parser    │  []Token → *ast.Program
 └──────┬───────┘
        ▼
 ┌──────────────┐  tabela de símbolos e tipos  semantico.md
 │   semantic   │  anota a AST com tipos
 └──────┬───────┘
        ▼
 ┌──────────────┐  dobra de constantes,        backend.md
 │  optimizer   │  remoção de código morto
 └──────┬───────┘
        ▼
 ┌──────────────┐  interpretação da AST        backend.md
 │ interpreter  │  stdin → stdout
 └──────────────┘
```

## Pacotes

| Pacote | Arquivos | Responsabilidade | API principal |
| --- | --- | --- | --- |
| `lexer` | `lexer.go`, `token.go` | Implementa o AFD; define os tipos de token e as palavras reservadas | `lexer.New(src).ScanTokens()`, `Errors()` |
| `ast` | `ast.go` | Nós da árvore; cada nó guarda o token de origem (linha e coluna) e as expressões guardam o tipo inferido | `ast.Program`, `ast.Expression.GetType()` |
| `parser` | `parser.go` | Descida recursiva sobre a GLC, constrói a AST, recupera erros | `parser.New(tokens).Parse()` → `(*ast.Program, []error)` |
| `semantic` | `analyzer.go`, `symbol_table.go` | Escopos, declarações, verificação e inferência de tipos | `semantic.New().Analyze(prog)`, `SymbolTable()` |
| `value` | `value.go` | Representação de valores, operadores, coerção e formatação; usado pelo otimizador e pelo interpretador | `Binary`, `Unary`, `Coerce`, `Zero`, `Format`, `Parse` |
| `optimizer` | `optimizer.go` | Dobra de constantes e eliminação de código morto sobre a AST | `optimizer.Optimize(prog)` → relatório |
| `interpreter` | `interpreter.go`, `environment.go` | Executa a AST com ambientes aninhados | `interpreter.New(in, out).Run(prog)` |
| `cmd/minilang` | `main.go` | CLI: escolhe o modo e encadeia as fases | `minilang [parse\|check\|opt\|run] <arquivo>` |

```text
minilang/
├── go.mod
├── cmd/minilang/        CLI e testes de integração (testdata/run/)
├── lexer/               M1
├── ast/                 M2
├── parser/              M2 (testdata/valid/, testdata/invalid/)
├── semantic/            M3
├── value/               M4
├── optimizer/           M4
├── interpreter/         M4
└── examples/            programas de exemplo
```

## Decisões que atravessam as fases

**Posição em todo nó.** Cada nó da AST carrega o `lexer.Token` que o originou.
Assim, as fases semântica e de execução conseguem relatar linha e coluna sem
guardar nenhuma estrutura extra.

**Erros acumulados.** Lexer, parser e analisador semântico devolvem uma lista
de erros e continuam a análise. O parser usa `panic`/`recover` internamente e
sincroniza em `;`, `}` ou no início do próximo comando. O analisador semântico
marca expressões com erro com o tipo `unknown`, o que evita erros em cascata.
A CLI só avança para a próxima fase se a lista estiver vazia.

**Um único modelo de valores.** Valores em tempo de execução são `int64`,
`float64` ou `bool`, e o otimizador e o interpretador fazem as contas pelas
mesmas funções do pacote `value`. Por isso, uma expressão dobrada em tempo de
compilação dá exatamente o mesmo resultado que daria na execução.

**Extensão C em todas as fases.** O tipo `real` aparece como `FLOAT_LITERAL` e
`TYPE_FLOAT` no lexer, como `real_literal` na gramática, como promoção
`inteiro → real` no analisador e como `float64` no pacote `value`.

## Formato dos erros

Todas as fases usam o mesmo formato:

```text
[ERRO <FASE>] linha <L>, coluna <C>: <descrição>
```

`<FASE>` é `LÉXICO`, `SINTÁTICO`, `SEMÂNTICO` ou `DE EXECUÇÃO`. Exemplos reais
estão no [guia de uso](../projeto/guia-de-uso.md#mensagens-de-erro).

## Testes

```bash
cd minilang
go test ./...
```

| Pacote | Tipo de teste | Onde ficam os casos |
| --- | --- | --- |
| `lexer` | Unitário: literal inteiro, comentário, `//`, `!` isolado | `lexer_test.go` |
| `parser` | Unitário (precedência, associatividade, dangling else, recuperação) e bateria de arquivos aceitos e rejeitados | `parser_test.go`, `parser/testdata/valid/*.min`, `parser/testdata/invalid/*.min` + `.want` |
| `semantic` | Unitário: escopos, tipos, coerção | `analyzer_test.go` |
| `value` | Unitário: operações e formatação | `value_test.go` |
| `optimizer` | Unitário: dobra e código morto | `optimizer_test.go` |
| `interpreter` | Unitário: execução de comandos | `interpreter_test.go` |
| `cmd/minilang` | Integração: CLI de ponta a ponta | `cmd/minilang/testdata/run/` |

Para adicionar um caso de integração, crie em `cmd/minilang/testdata/run/`:

| Arquivo | Conteúdo |
| --- | --- |
| `nome.min` | Programa |
| `nome.want` | Saída padrão esperada |
| `nome.in` | Entrada padrão (opcional) |
| `nome.err` | Trecho esperado da saída de erro; indica que o programa deve falhar (opcional) |

Para o parser, coloque o arquivo em `parser/testdata/valid/` ou em
`parser/testdata/invalid/` (com um `.want` contendo os diagnósticos). Os testes
leem esses diretórios automaticamente. A bateria oficial da disciplina pode
ser adicionada da mesma forma.

## Limitações conhecidas

* O lexer escreve os erros léxicos diretamente na saída padrão
  (`fmt.Printf` em `reportInvalidCharacter`), além de devolvê-los em
  `Errors()`. Por isso aparecem no meio da listagem de tokens e não passam pelo
  `errorOutput` da CLI.
* Identificadores só aceitam letras ASCII, `ã` e `Ã`.
* A execução para no primeiro erro de execução.

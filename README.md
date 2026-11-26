# MiniLang Compiler

Compilador/interpretador desenvolvido em **Go** para a linguagem **MiniLang**, como parte da Avaliação A3 da disciplina **Teoria da Computação e Compiladores — UNIFACS 2026.2**.

O projeto será desenvolvido de forma incremental através dos marcos **M1, M2, M3 e M4**, contemplando as principais etapas de construção de um compilador: análise léxica, análise sintática, construção da AST, análise semântica e back-end.

## 📌 Sobre o projeto

A **MiniLang** é uma linguagem imperativa de propósito educacional cujo escopo prevê:

* Declaração de variáveis;
* Tipos `inteiro` e `booleano`;
* Atribuições;
* Entrada e saída;
* Condicionais `se` e `senão`;
* Laços `enquanto`;
* Blocos de comandos;
* Expressões aritméticas;
* Expressões relacionais;
* Expressões lógicas;
* Comentários;
* Detecção de erros léxicos, sintáticos e semânticos.

Como extensão do projeto, será implementado também o tipo **`real`**, incluindo coerção automática de valores `inteiro` para `real`.

---

## 🚀 Tecnologias

O compilador será desenvolvido utilizando:

* **Go**
* Testes nativos do Go (`go test`)

A implementação dos principais componentes do compilador será feita manualmente, permitindo explorar diretamente os conceitos de autômatos, análise léxica, gramáticas, parsing, AST, análise semântica e interpretação.

---

## 🔄 Pipeline do compilador

O processamento de um programa MiniLang seguirá, ao final do projeto, o seguinte fluxo:

```text
Código-fonte (.min)
        │
        ▼
┌───────────────────┐
│ Analisador Léxico │
│       M1          │
└─────────┬─────────┘
          │
        Tokens
          │
          ▼
┌───────────────────────┐
│ Analisador Sintático  │
│         M2            │
└──────────┬────────────┘
           │
          AST
           │
           ▼
┌───────────────────────┐
│ Analisador Semântico  │
│         M3            │
└──────────┬────────────┘
           │
      AST validada
           │
           ▼
┌───────────────────────┐
│      Back-end         │
│         M4            │
└──────────┬────────────┘
           │
           ▼
        Execução
```

---

# 🗺️ Marcos do projeto

## M1 — Analisador Léxico

O primeiro marco consiste na implementação do analisador léxico da MiniLang.

O lexer já transforma o código-fonte em uma sequência de **tokens**.

Exemplo:

```text
var idade: inteiro;
idade = 20;
```

Resultado conceitual esperado:

```text
VAR
IDENTIFIER("idade")
COLON
INTEGER_TYPE
SEMICOLON

IDENTIFIER("idade")
ASSIGN
INTEGER_LITERAL("20")
SEMICOLON
```

### Requisitos do M1

* Reconhecimento das palavras reservadas;
* Reconhecimento de identificadores;
* Reconhecimento de números;
* Reconhecimento de operadores;
* Reconhecimento de delimitadores;
* Tratamento de `=`, `==`, `<`, `<=`, `>`, `>=` e `!=`;
* Descarte de espaços em branco;
* Descarte de comentários;
* Rastreamento de linha e coluna;
* Detecção de caracteres inválidos;
* Mensagens de erro léxico;
* Documentação do AFD em [`docs/afd.md`](docs/afd.md);
* Testes automatizados.

---

## M2 — Analisador Sintático + AST

O segundo marco adicionará o **parser**.

O parser receberá os tokens produzidos pelo lexer e verificará se eles seguem a gramática livre de contexto definida para a MiniLang. A gramática EBNF e a descrição do autômato de pilha correspondente estão em [`docs/grammar.ebnf`](docs/grammar.ebnf) e [`docs/pda.md`](docs/pda.md).

```text
Código
  ↓
Lexer
  ↓
Tokens
  ↓
Parser
  ↓
AST
```

Além da validação sintática, será construída uma **Abstract Syntax Tree (AST)**.

Exemplo:

```text
resultado = 10 + 5 * 2;
```

AST conceitual:

```text
Assignment
├── Identifier: resultado
└── Addition
    ├── IntegerLiteral: 10
    └── Multiplication
        ├── IntegerLiteral: 5
        └── IntegerLiteral: 2
```

O parser deverá respeitar corretamente precedência e associatividade dos operadores.

O parser descendente recursivo constrói a AST, aplica precedência e associatividade e recupera erros em pontos de sincronização para permitir a identificação de múltiplos erros em uma única execução. Os vocábulos em português e os aliases em inglês são reconhecidos pelo lexer.

Na API Go, `parser.New(tokens).Parse()` retorna a AST parcial ou completa e uma lista de erros sintáticos. Cada nó conserva o token de origem, incluindo sua linha e coluna.

---

## M3 — Analisador Semântico

O terceiro marco será responsável pela análise semântica do programa.

Será criada uma **tabela de símbolos** contendo informações como:

```text
Identificador | Tipo     | Escopo | Posição
--------------|----------|--------|--------
idade         | inteiro  | global | 2:5
altura        | real     | global | 3:5
ativo         | booleano | global | 4:5
```

A análise semântica deverá detectar situações como:

* Uso de variável não declarada;
* Redeclaração;
* Incompatibilidade de tipos;
* Condições não booleanas;
* Uso incorreto de operadores aritméticos;
* Uso incorreto de operadores lógicos;
* Uso incorreto do comando `leia`;
* Outros erros semânticos definidos pela linguagem.

A AST também será anotada com os tipos inferidos durante a análise.

No marco M3, os blocos criam escopos aninhados, declarações são visíveis a
partir do ponto em que aparecem e a extensão `real` permite promoção de
`inteiro` para `real`. A implementação e suas regras estão documentadas em
[`docs/semantic.md`](docs/semantic.md).

---

# 🔢 Extensão C — Tipo `real`

A extensão escolhida para o projeto é:

> **Extensão C — Tipo real com coerção de inteiro para real.**

A MiniLang originalmente possui os tipos:

```text
inteiro
booleano
```

Nossa implementação adicionará:

```text
real
```

Permitindo programas como:

```text
var preco: real;
var quantidade: inteiro;
var total: real;

preco = 19.90;
quantidade = 3;

total = preco * quantidade;

escreva(total);
```

---

## Promoção de tipos

A extensão permitirá a promoção automática:

```text
inteiro → real
```

Por exemplo:

```text
10 + 2.5
```

Os operandos possuem tipos diferentes:

```text
10  → inteiro
2.5 → real
```

Durante a análise semântica, o valor inteiro poderá ser promovido para `real`:

```text
inteiro + real
       ↓
 real + real
       ↓
     real
```

Portanto:

```text
10 + 2.5
```

resultará em:

```text
12.5
```

### Regras iniciais de operações numéricas

| Operando esquerdo | Operando direito | Resultado |
| ----------------- | ---------------- | --------- |
| `inteiro`         | `inteiro`        | `inteiro` |
| `inteiro`         | `real`           | `real`    |
| `real`            | `inteiro`        | `real`    |
| `real`            | `real`           | `real`    |

Operações aritméticas envolvendo valores booleanos deverão gerar erro semântico.

Exemplo inválido:

```text
verdadeiro + 10.5
```

---

## Atribuições e coerção

Será permitida atribuição de um `inteiro` para uma variável `real`:

```text
var valor: real;

valor = 10;
```

Conceitualmente:

```text
10
│
inteiro
│
▼
coerção
│
▼
10.0
│
real
```

A conversão implícita inversa não faz parte da extensão definida:

```text
var quantidade: inteiro;

quantidade = 5.75;
```

deverá resultar em erro semântico.

Exemplo:

```text
Erro semântico [linha 3, coluna 14]:
não é possível atribuir um valor do tipo real
a uma variável do tipo inteiro.
```

---

# ⚙️ M4 — Back-end

No último marco será implementado o back-end da linguagem.

A estratégia inicialmente planejada é a **interpretação direta da AST**.

Exemplo:

```text
var a: inteiro;
var b: real;
var resultado: real;

a = 10;
b = 2.5;

resultado = a + b;

escreva(resultado);
```

Pipeline:

```text
Código
   ↓
Lexer
   ↓
Tokens
   ↓
Parser
   ↓
AST
   ↓
Análise Semântica
   ↓
AST validada
   ↓
Interpretador
   ↓
12.5
```

O back-end contempla a otimização exigida pelo projeto: **dobra de constantes**
(`2 * (3 + 4) - 1` → `13`) e **eliminação de código morto** (`se` com condição
constante e `enquanto (falso)`). Interpretador, regras de execução, erros de
execução e otimização estão documentados em [`docs/backend.md`](docs/backend.md).

---

# 📁 Estrutura atual e planejada

Atualmente, a raiz contém este README, a licença e a documentação do AFD em [`docs/afd.md`](docs/afd.md). O módulo Go fica em `minilang/`, com `go.mod`, `cmd/minilang/main.go`, `lexer/lexer.go`, `lexer/token.go` e os exemplos `examples/exemple.min` e `examples/exemple-identifer.min`.

A árvore abaixo resume a organização do módulo. Lexer, parser, AST, análise
semântica, otimizador e interpretador estão implementados, com seus testes.
Este README e o diretório `docs/` ficam na raiz do repositório.

```text
minilang/
│
├── cmd/
│   └── minilang/
│       ├── main.go
│       ├── main_test.go
│       └── testdata/run/
│
├── lexer/
│   ├── lexer.go
│   ├── token.go
│   └── lexer_test.go
│
├── parser/
│   ├── parser.go
│   └── parser_test.go
│
├── ast/
│   └── ast.go
│
├── semantic/
│   ├── analyzer.go
│   ├── symbol_table.go
│   └── analyzer_test.go
│
├── value/
│   ├── value.go
│   └── value_test.go
│
├── optimizer/
│   ├── optimizer.go
│   └── optimizer_test.go
│
├── interpreter/
│   ├── interpreter.go
│   ├── environment.go
│   └── interpreter_test.go
│
├── examples/
│   ├── parser/
│   ├── programa.min
│   └── real.min
│
└── go.mod
```

```text
docs/
├── afd.md
├── grammar.ebnf
├── pda.md
├── semantic.md
└── backend.md
```

A estrutura poderá sofrer alterações conforme a evolução do projeto.

---

# 💻 Executando o projeto

Pré-requisito: Go instalado e disponível no `PATH`. O `minilang/go.mod` declara a versão **1.27.1**. A partir da raiz do repositório, entre no módulo:

```bash
cd minilang
go version
```

Executar a análise léxica dos exemplos existentes:

```bash
go run ./cmd/minilang examples/exemple.min
go run ./cmd/minilang examples/exemple-identifer.min
```

O primeiro exemplo exercita operadores e delimitadores. O segundo inclui identificadores, números e um erro proposital com `@`. Comentários MiniLang começam com `#`; a sequência `//` representa dois operadores de divisão.

Executar a análise sintática dos exemplos M2:

```bash
go run ./cmd/minilang parse examples/parser/valido.min
go run ./cmd/minilang parse examples/parser/erro-declaracao.min
go run ./cmd/minilang parse examples/parser/multiplos-erros.min
```

Executar a análise semântica (lexer + parser + semântico):

```bash
go run ./cmd/minilang check examples/parser/valido.min
```

O modo `check` informa erros semânticos com linha e coluna ou, em caso de
sucesso, imprime a tabela de símbolos com tipo, escopo e posição de declaração.

Executar um programa (pipeline completo, com otimização e interpretador):

```bash
echo 6 | go run ./cmd/minilang run examples/programa.min
go run ./cmd/minilang run examples/real.min
```

`leia` consome valores da entrada padrão. Erros de qualquer fase interrompem a
execução com código de saída diferente de zero. Para ver o efeito da
otimização:

```bash
go run ./cmd/minilang opt examples/programa.min
```

O modo `parse` mostra os diagnósticos sintáticos com linha e coluna. Para programas aceitos, confirma a construção da AST. Para programas inválidos, retorna código de saída diferente de zero; o segundo exemplo exibe dois erros na mesma execução. A execução sem `parse` continua apresentando os tokens do lexer.

A saída apresenta o tipo do token, o lexema entre aspas e a posição `linha:coluna`. Por exemplo, `20.5;` gera `FLOAT_LITERAL` e `SEMICOLON`. Ao final do arquivo é emitido `EOF`. Sem um caminho de arquivo, a CLI exibe `Uso: minilang <arquivo>`.

---

# 🧪 Testes

Os testes automatizados usam o pacote nativo de testes do Go. A bateria de programas do parser fica em `minilang/parser/testdata/`: `valid/` contém fontes aceitas e `invalid/` contém fontes rejeitadas, com diagnósticos esperados nos arquivos `.want`. Os testes leem esses casos automaticamente. Execute os comandos abaixo dentro de `minilang/`.

Executar todos os testes:

```bash
go test ./...
```

Executar especificamente os testes do lexer:

```bash
go test ./lexer
```

Executar a bateria do parser com cada caso listado:

```bash
go test -v ./parser
```

Executar com informações detalhadas:

```bash
go test -v ./...
```

Essa bateria local é reproduzível e cobre casos válidos, mensagens e posições dos erros, além de múltiplos erros em uma execução. Quando a bateria oficial da disciplina estiver disponível, seus arquivos podem ser adicionados a `testdata/` para validar a cobertura completa exigida pela avaliação.

---

# ❌ Tratamento de erros

Caracteres inválidos, incluindo `@` e `!` isolado, geram um diagnóstico com linha e coluna, e a análise continua até `EOF`. As mensagens identificam:

* Fase do compilador;
* Linha;
* Coluna;
* Descrição do problema.

Exemplo léxico:

```text
[ERRO LÉXICO] linha 4, coluna 10:
caractere '@' não reconhecido.
```

Exemplo sintático:

```text
[ERRO SINTÁTICO] linha 7, coluna 15:
esperado ';' após atribuição.
```

Exemplo semântico:

```text
[ERRO SEMÂNTICO] linha 10, coluna 5:
variável 'resultado' não declarada.
```

Exemplo relacionado à extensão C:

```text
[ERRO SEMÂNTICO] linha 12, coluna 14:
não é possível atribuir REAL a uma variável INTEIRO.
```

Exemplo de execução:

```text
[ERRO DE EXECUÇÃO] linha 5, coluna 11:
divisão por zero.
```

---

# 📝 Palavras reservadas

O mapa `Keywords` em `minilang/lexer/token.go` reconhece as palavras abaixo, diferenciando maiúsculas de minúsculas. Cada forma em português e seu alias em inglês geram o mesmo tipo de token:

```text
program/programa var int/inteiro float/real boolean/booleano
if/se else/senão while/enquanto print/escreva read/leia
true/verdadeiro false/falso and/e or/ou not/não end/fim
```

Os dois vocabulários são aceitos para manter compatibilidade com os exemplos já existentes e com a especificação da linguagem.

---

# 🔣 Operadores

### Aritméticos

```text
+  -  *  /  %
```

### Relacionais

```text
==  !=  <  <=  >  >=
```

### Lógicos

O lexer atual reconhece `and`, `or` e `not`. Na especificação em português:

```text
e
ou
não
```

### Atribuição

```text
=
```

---

# 💬 Comentários

Comentários começam com `#` e continuam até o final da linha.

```text
# Isto é um comentário

var idade: inteiro; # declaração de idade
```

Os comentários são descartados durante a análise léxica. A sequência `//` não é comentário e gera dois tokens `DIVIDE`.

---

# 🗓️ Roadmap

### M1 — Analisador Léxico

* [x] Estrutura inicial do projeto Go e CLI
* [x] Definição dos tokens
* [x] Palavras reservadas em inglês
* [ ] Alinhar palavras reservadas com a especificação em português
* [x] Identificadores
* [x] Literais inteiros (`INT_LITERAL`)
* [x] Literais reais
* [x] Operadores
* [x] Delimitadores
* [x] Comentários com `#`
* [x] Linha e coluna
* [x] Erros léxicos (incluindo `!` isolado)
* [x] AFD documentado em [`docs/afd.md`](docs/afd.md)
* [x] Testes do lexer para os comportamentos implementados

### M2 — Parser + AST

* [x] Gramática EBNF e documentação da GLC/PDA
* [x] Parser descendente recursivo
* [x] AST de comandos e expressões com localização dos tokens
* [x] Precedência e associatividade
* [x] Literais reais na AST
* [x] Recuperação de erros
* [x] Dangling else
* [x] Testes sintáticos

### M3 — Análise Semântica

* [x] Tabela de símbolos
* [x] Escopos
* [x] Verificação de declarações
* [x] Sistema de tipos
* [x] Tipo `real`
* [x] Promoção `inteiro → real`
* [x] Verificação de operadores
* [x] AST anotada
* [x] Testes semânticos

### M4 — Back-end

* [x] Interpretador da AST
* [x] Execução de valores inteiros
* [x] Execução de valores reais
* [x] Coerção `inteiro → real`
* [x] Entrada e saída
* [x] Condicionais
* [x] Laços
* [x] Otimização
* [x] Testes de integração
* [ ] Relatório técnico
* [ ] Preparação da demonstração

---

# 🎯 Objetivo final

Dentro de `minilang/`, um programa completo como `examples/programa.min` é executado com:

```bash
go run ./cmd/minilang run examples/programa.min
```

que realiza automaticamente:

```text
programa.min
     ↓
   Lexer
     ↓
   Tokens
     ↓
   Parser
     ↓
    AST
     ↓
Semântico
     ↓
Interpretador
     ↓
  Resultado
```

O objetivo é que o projeto demonstre de forma prática os conceitos estudados em **Teoria da Computação e Compiladores**, relacionando modelos formais de reconhecimento de linguagens à implementação de um compilador/interpretador funcional.

---

## 👥 Equipe

**Instituição:** UNIFACS
**Disciplina:** Teoria da Computação e Compiladores
**Período:** 2026.2

### Integrantes

* Hanspeter Dietiker
* Gabriel Moraes


---

## 📚 Referências

* AHO, Alfred V. et al. *Compiladores: princípios, técnicas e ferramentas*. 2ª ed. Pearson Addison Wesley, 2008.
* MENEZES, Paulo Blauth. *Linguagens Formais e Autômatos*. 6ª ed. Bookman, 2011.
* SANTOS, Pedro Reis; LANGLOIS, Thibault. *Compiladores: da teoria à prática*. LTC, 2018.
* Material da disciplina Teoria da Computação e Compiladores — UNIFACS 2026.2.

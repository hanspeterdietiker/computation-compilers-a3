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

O parser receberá os tokens produzidos pelo lexer e verificará se eles seguem a gramática definida para a MiniLang.

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

Também serão implementadas estratégias de recuperação de erros sintáticos para permitir a identificação de múltiplos erros em uma única execução.

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

O back-end também deverá contemplar a otimização exigida pelo projeto.

---

# 📁 Estrutura atual e planejada

Atualmente, a raiz contém este README, a licença e a documentação do AFD em [`docs/afd.md`](docs/afd.md). O módulo Go fica em `minilang/`, com `go.mod`, `cmd/minilang/main.go`, `lexer/lexer.go`, `lexer/token.go` e os exemplos `examples/exemple.min` e `examples/exemple-identifer.min`.

A árvore abaixo representa a organização planejada; arquivos de testes, parser, AST, semântica e interpretador ainda não existem. Este README e o diretório `docs/` ficam na raiz do repositório.

```text
minilang/
│
├── cmd/
│   └── minilang/
│       └── main.go
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
├── interpreter/
│   ├── interpreter.go
│   └── interpreter_test.go
│
├── examples/
│   ├── valido.min
│   ├── invalido.min
│   └── real.min
│
└── go.mod
```

```text
docs/
├── afd.md
└── grammar.ebnf  # planejado para o M2
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

O primeiro exemplo exercita operadores e delimitadores. O segundo inclui identificadores, números e um erro proposital com `@`. Comentários MiniLang começam com `#`; a sequência `//` representa dois operadores de divisão. No segundo exemplo, `inteiro` é reconhecido como identificador; a palavra reservada implementada é `int`.

A saída apresenta o tipo do token, o lexema entre aspas e a posição `linha:coluna`. Por exemplo, `20.5;` gera `FLOAT_LITERAL` e `SEMICOLON`. Ao final do arquivo é emitido `EOF`. Sem um caminho de arquivo, a CLI exibe `Uso: minilang <arquivo>`.

---

# 🧪 Testes

Os testes automatizados usam o pacote nativo de testes do Go. A suíte do lexer cobre literais inteiros, comentários com `#`, o significado de `//` e o erro para `!` isolado. Execute os comandos abaixo dentro de `minilang/`.

Executar todos os testes:

```bash
go test ./...
```

Executar especificamente os testes do lexer:

```bash
go test ./lexer
```

Executar com informações detalhadas:

```bash
go test -v ./...
```

A bateria de testes deverá possuir tanto programas válidos quanto inválidos.

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

---

# 📝 Palavras reservadas

O mapa `Keywords` em `minilang/lexer/token.go` reconhece atualmente as seguintes palavras, diferenciando maiúsculas de minúsculas:

```text
program var int float boolean
if else while print read
true false and or not end
```

O alinhamento com a especificação em português permanece pendente. As palavras previstas nessa especificação são:

```text
programa
var
inteiro
booleano
se
senão
enquanto
escreva
leia
verdadeiro
falso
e
ou
não
fim
```

Com a Extensão C adicionaremos:

```text
real
```

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

* [ ] Gramática EBNF
* [ ] Parser
* [ ] AST
* [ ] Precedência
* [ ] Associatividade
* [ ] Literais reais na AST
* [ ] Recuperação de erros
* [ ] Dangling else
* [ ] Testes sintáticos

### M3 — Análise Semântica

* [ ] Tabela de símbolos
* [ ] Escopos
* [ ] Verificação de declarações
* [ ] Sistema de tipos
* [ ] Tipo `real`
* [ ] Promoção `inteiro → real`
* [ ] Verificação de operadores
* [ ] AST anotada
* [ ] Testes semânticos

### M4 — Back-end

* [ ] Interpretador da AST
* [ ] Execução de valores inteiros
* [ ] Execução de valores reais
* [ ] Coerção `inteiro → real`
* [ ] Entrada e saída
* [ ] Condicionais
* [ ] Laços
* [ ] Otimização
* [ ] Testes de integração
* [ ] Relatório técnico
* [ ] Preparação da demonstração

---

# 🎯 Objetivo final

Ao final do projeto, esperamos executar, dentro de `minilang/`, um programa completo como o futuro `examples/programa.min`:

```bash
go run ./cmd/minilang examples/programa.min
```

e realizar automaticamente:

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

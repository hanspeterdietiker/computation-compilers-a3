# MiniLang Compiler

Compilador/interpretador desenvolvido em **Go** para a linguagem **MiniLang**, como parte da Avaliação A3 da disciplina **Teoria da Computação e Compiladores — UNIFACS 2026.2**.

O projeto será desenvolvido de forma incremental através dos marcos **M1, M2, M3 e M4**, contemplando as principais etapas de construção de um compilador: análise léxica, análise sintática, construção da AST, análise semântica e back-end.

## 📌 Sobre o projeto

A **MiniLang** é uma linguagem imperativa de propósito educacional que possui suporte a:

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

O lexer será responsável por transformar o código-fonte em uma sequência de **tokens**.

Exemplo:

```text
var idade: inteiro;
idade = 20;
```

Resultado conceitual:

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
* Documentação do AFD;
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

# 📁 Estrutura planejada

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
├── docs/
│   ├── afd.md
│   └── grammar.ebnf
│
├── go.mod
└── README.md
```

A estrutura poderá sofrer alterações conforme a evolução do projeto.

---

# 💻 Executando o projeto

> Esta seção será atualizada conforme a implementação avançar.

Pré-requisito:

```bash
go version
```

Executar o compilador:

```bash
go run ./cmd/minilang examples/valido.min
```

---

# 🧪 Testes

Os testes automatizados serão implementados utilizando o pacote nativo de testes do Go.

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

As mensagens de erro deverão identificar:

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

A especificação mínima possui as seguintes palavras reservadas:

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

Os comentários serão descartados durante a análise léxica.

---

# 🗓️ Roadmap

### M1 — Analisador Léxico

* [ ] Estrutura inicial do projeto Go
* [ ] Definição dos tokens
* [ ] Palavras reservadas
* [ ] Identificadores
* [ ] Literais inteiros
* [ ] Literais reais
* [ ] Operadores
* [ ] Delimitadores
* [ ] Comentários
* [ ] Linha e coluna
* [ ] Erros léxicos
* [ ] AFD
* [ ] Testes do lexer

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

Ao final do projeto, esperamos executar:

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

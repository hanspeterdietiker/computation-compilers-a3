# MiniLang

Compilador/interpretador da linguagem educacional **MiniLang**, escrito em **Go**
para a Avaliação A3 de **Teoria da Computação e Compiladores (UNIFACS, 2026.2)**.

Ele lê um programa `.min`, faz as análises léxica, sintática e semântica, otimiza
a árvore sintática e executa o programa. Além da linguagem base, implementamos a
**Extensão C: tipo `real` com coerção automática de `inteiro` para `real`**.

```text
var preco: real;
var quantidade: inteiro;
var total: real;

preco = 19.90;
quantidade = 3;
total = preco * quantidade;   # inteiro promovido para real
escreva(total);               # 59.7
```

## Início rápido

Requer Go **1.27.1** ou superior.

```bash
cd minilang
go run ./cmd/minilang run examples/real.min         # executa um programa
echo 6 | go run ./cmd/minilang run examples/programa.min
go test ./...                                       # roda todos os testes
```

## Documentação

| Para quem | Documento | Conteúdo |
| --- | --- | --- |
| Quem quer entender o projeto | [Visão geral](docs/projeto/visao-geral.md) | O problema, o objetivo, a extensão escolhida, o que está pronto |
| Quem vai escrever programas | [Linguagem MiniLang](docs/projeto/linguagem.md) | Sintaxe, tipos, operadores e exemplos |
| Quem vai usar a ferramenta | [Guia de uso](docs/projeto/guia-de-uso.md) | Comandos da CLI, saídas e mensagens de erro |
| Quem vai ler ou alterar o código | [Documentação técnica](docs/tecnico/arquitetura.md) | Arquitetura, AFD, gramática, AP, semântica e back-end |

O índice completo está em [`docs/README.md`](docs/README.md).

## Equipe

* Hanspeter Dietiker
* Gabriel Moraes

Licença: ver [`LICENSE`](LICENSE).

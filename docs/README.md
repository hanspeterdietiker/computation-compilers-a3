# Documentação da MiniLang

A documentação tem duas partes. Leia a do **projeto** para saber o que o
compilador faz e como usá-lo. Leia a **técnica** para entender como ele foi
construído.

## Projeto: o quê e para quê

| Documento | Pergunta que responde |
| --- | --- |
| [Visão geral](projeto/visao-geral.md) | Qual problema resolvemos, por que e o que está pronto? |
| [Linguagem MiniLang](projeto/linguagem.md) | Como se escreve um programa MiniLang? |
| [Guia de uso](projeto/guia-de-uso.md) | Como rodar o compilador e ler as mensagens de erro? |

## Técnica: como funciona

Os documentos seguem a ordem do pipeline. Cada etapa recebe o resultado da
anterior.

| Documento | Etapa | Modelo formal | Pacote Go |
| --- | --- | --- | --- |
| [Arquitetura](tecnico/arquitetura.md) | Visão do pipeline, pacotes, testes | — | todos |
| [Análise léxica](tecnico/lexico.md) | Código-fonte → tokens | AFD | `lexer` |
| [Gramática](tecnico/gramatica.ebnf) | Definição da sintaxe | GLC (EBNF) | — |
| [Análise sintática](tecnico/sintatico.md) | Tokens → AST | Autômato de pilha | `parser`, `ast` |
| [Análise semântica](tecnico/semantico.md) | AST → AST anotada | Tabela de símbolos e tipos | `semantic` |
| [Back-end](tecnico/backend.md) | AST anotada → execução | Otimização e interpretação | `value`, `optimizer`, `interpreter` |

```text
docs/
├── README.md              este índice
├── projeto/
│   ├── visao-geral.md
│   ├── linguagem.md
│   └── guia-de-uso.md
└── tecnico/
    ├── arquitetura.md
    ├── lexico.md
    ├── gramatica.ebnf
    ├── sintatico.md
    ├── semantico.md
    └── backend.md
```

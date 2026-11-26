# Visão geral do projeto

## O problema

A Avaliação A3 de **Teoria da Computação e Compiladores (UNIFACS, 2026.2)**
pede a construção de um compilador completo para a **MiniLang**, uma linguagem
imperativa pequena criada para fins didáticos. O trabalho precisa mostrar, em
código funcionando, os modelos formais vistos na disciplina:

* um **autômato finito determinístico (AFD)** que reconhece os tokens;
* uma **gramática livre de contexto (GLC)** e o **autômato de pilha** que
  reconhece a sintaxe;
* as **regras sensíveis ao contexto** (declarações, escopos e tipos) que a
  gramática sozinha não consegue expressar;
* um **back-end** que otimiza e executa o programa.

A linguagem base tem os tipos `inteiro` e `booleano`, declarações, atribuições,
`leia`/`escreva`, `se`/`senão`, `enquanto`, blocos, expressões aritméticas,
relacionais e lógicas e comentários. Cada equipe escolhe também uma
**extensão** da linguagem.

## A extensão escolhida: Extensão C, tipo `real`

Escolhemos a **Extensão C: tipo `real` com coerção de `inteiro` para `real`**.
Ela acrescenta números com casas decimais e obriga o compilador a tratar dois
tipos numéricos ao mesmo tempo. Isso afeta todas as etapas:

| Etapa | O que a extensão exige |
| --- | --- |
| Léxica | Reconhecer literais como `19.90` (`FLOAT_LITERAL`) sem confundir com `19` seguido de `.` |
| Sintática | Aceitar `real` como tipo e literais reais em expressões |
| Semântica | Promover `inteiro` para `real` em operações mistas e em atribuições; proibir `real → inteiro` |
| Back-end | Calcular com `float64`, converter valores e imprimir reais com parte decimal |

As regras da extensão:

| Operando esquerdo | Operando direito | Resultado de `+ - * /` |
| --- | --- | --- |
| `inteiro` | `inteiro` | `inteiro` |
| `inteiro` | `real` | `real` |
| `real` | `inteiro` | `real` |
| `real` | `real` | `real` |

```text
var valor: real;
valor = 10;          # permitido: 10 vira 10.0

var quantidade: inteiro;
quantidade = 5.75;   # erro semântico: não há conversão implícita real → inteiro
```

Operações aritméticas com `booleano` continuam proibidas
(`verdadeiro + 10.5` é erro semântico).

## O que o projeto faz

O programa `minilang` recebe um arquivo `.min` e:

1. **separa o texto em tokens** e aponta caracteres inválidos com linha e coluna;
2. **verifica a sintaxe**, monta a árvore sintática abstrata (AST) e informa
   vários erros sintáticos numa mesma execução;
3. **verifica o significado**: variáveis declaradas, escopos, tipos e coerções;
4. **otimiza** o programa (dobra de constantes e remoção de código morto);
5. **executa** o programa, lendo da entrada padrão e escrevendo na saída padrão.

Cada etapa pode ser executada isoladamente pela CLI (`tokens`, `parse`, `check`,
`opt`, `run`), o que permite mostrar o resultado intermediário de cada fase
durante a apresentação. Veja o [guia de uso](guia-de-uso.md).

## Para que serve

* **Para a avaliação:** cada marco (M1 a M4) corresponde a uma etapa entregue,
  testada e documentada.
* **Para estudo:** o código foi escrito à mão, sem geradores de lexer ou de
  parser, para que a ligação entre o modelo formal e a implementação fique
  visível. Cada documento técnico parte do modelo (AFD, GLC, AP) e mostra a
  função Go correspondente.
* **Para experimentação:** é possível escrever e executar programas MiniLang
  reais, como o cálculo de soma e média em
  [`examples/programa.min`](../../minilang/examples/programa.min).

## Decisões de projeto

| Decisão | Motivo |
| --- | --- |
| Go, sem bibliotecas externas | Binário único, testes nativos (`go test`), código fácil de ler |
| Lexer e parser escritos à mão | Mostrar o AFD e a descida recursiva diretamente no código |
| Interpretação direta da AST | Back-end mais simples, adequado ao escopo; a otimização atua na própria AST |
| Palavras reservadas em português e em inglês | `se`/`if`, `enquanto`/`while` etc. geram o mesmo token; mantém compatibilidade com os primeiros exemplos |
| Recuperação de erros | Relatar todos os erros de uma vez em vez de parar no primeiro |

## Situação dos marcos

| Marco | Entrega | Situação |
| --- | --- | --- |
| M1 | Analisador léxico, AFD documentado, testes | Concluído |
| M2 | Gramática EBNF, autômato de pilha, parser, AST, recuperação de erros | Concluído |
| M3 | Tabela de símbolos, escopos, tipos, extensão `real`, AST anotada | Concluído |
| M4 | Otimizador, interpretador, testes de integração | Concluído |

Pendências:

* [ ] Relatório técnico final
* [ ] Preparação da demonstração
* [ ] Executar a bateria oficial da disciplina quando ela for disponibilizada
  (os arquivos podem ser copiados para os diretórios `testdata/`)

## Equipe

**Instituição:** UNIFACS. **Disciplina:** Teoria da Computação e Compiladores.
**Período:** 2026.2.

* Hanspeter Dietiker
* Gabriel Moraes

## Referências

* AHO, Alfred V. et al. *Compiladores: princípios, técnicas e ferramentas*. 2ª ed. Pearson Addison Wesley, 2008.
* MENEZES, Paulo Blauth. *Linguagens Formais e Autômatos*. 6ª ed. Bookman, 2011.
* SANTOS, Pedro Reis; LANGLOIS, Thibault. *Compiladores: da teoria à prática*. LTC, 2018.
* Material da disciplina Teoria da Computação e Compiladores, UNIFACS 2026.2.

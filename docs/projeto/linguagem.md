# Linguagem MiniLang

Guia para quem vai escrever programas. A definição formal da sintaxe está em
[`gramatica.ebnf`](../tecnico/gramatica.ebnf).

## Estrutura de um programa

O cabeçalho e o `fim` são opcionais. As três formas abaixo são válidas:

```text
programa demonstracao {
    escreva(1);
} fim
```

```text
programa exemplo
escreva(1);
fim
```

```text
escreva(1);
```

Arquivos MiniLang usam a extensão `.min`.

## Tipos

| Tipo | Valores | Valor inicial |
| --- | --- | --- |
| `inteiro` | `0`, `42`, `-7` | `0` |
| `real` | `2.5`, `19.90`, `10.0` | `0.0` |
| `booleano` | `verdadeiro`, `falso` | `falso` |

O tipo `real` é a extensão escolhida pela equipe. Um literal real precisa ter
dígitos dos dois lados do ponto: `2.5` é válido, `2.` e `.5` não são.

## Comandos

| Comando | Exemplo |
| --- | --- |
| Declaração | `var idade: inteiro;` |
| Atribuição | `idade = 20;` |
| Saída | `escreva(idade + 1);` |
| Entrada | `leia(idade);` |
| Condicional | `se (idade >= 18) escreva(1); senão escreva(0);` |
| Laço | `enquanto (i < 10) { i = i + 1; }` |
| Bloco | `{ var x: inteiro; x = 1; }` |

* Toda variável precisa ser declarada **antes** de ser usada.
* Cada bloco `{ ... }` abre um novo escopo. Uma variável declarada dentro do
  bloco não existe fora dele e pode esconder outra de mesmo nome declarada fora.
* Declarar duas vezes o mesmo nome no mesmo escopo é erro.
* As condições de `se` e `enquanto` precisam ser `booleano`.
* `senão` pertence ao `se` mais próximo que ainda não tem `senão`.
* `leia` lê um valor da entrada padrão e converte para o tipo da variável.
* `escreva` imprime um valor por linha. Reais sempre aparecem com parte decimal
  (`10.0`).

## Operadores

Da **menor** para a **maior** precedência:

| Precedência | Operadores | Operandos | Resultado |
| --- | --- | --- | --- |
| 1 | `ou` | `booleano` | `booleano` |
| 2 | `e` | `booleano` | `booleano` |
| 3 | `==` `!=` | dois numéricos ou dois `booleano` | `booleano` |
| 4 | `<` `<=` `>` `>=` | numéricos | `booleano` |
| 5 | `+` `-` | numéricos | `real` se algum for `real`, senão `inteiro` |
| 6 | `*` `/` `%` | numéricos (`%` só `inteiro`) | idem |
| 7 | `não`, `-` e `+` unários | conforme o operador | conforme o operador |

* Operadores binários associam à esquerda: `a - b - c` é `(a - b) - c`.
* `/` entre dois inteiros faz divisão inteira: `7 / 2` é `3`.
* `e` e `ou` param assim que o resultado é conhecido (curto-circuito).
* Parênteses mudam a ordem: `2 * (3 + 4)` é `14`.

## Coerção entre `inteiro` e `real`

O compilador converte `inteiro` para `real` automaticamente quando necessário.
O contrário nunca acontece de forma implícita.

```text
var r: real;
var i: inteiro;

r = 10;        # ok: r recebe 10.0
r = i * 2.5;   # ok: i é promovido para real
i = 2.5;       # erro: não é possível atribuir real a inteiro
```

## Comentários

Começam com `#` e vão até o fim da linha.

```text
# comentário de linha inteira
var idade: inteiro; # comentário no fim da linha
```

`//` não é comentário: são dois operadores de divisão.

## Palavras reservadas

As palavras diferenciam maiúsculas de minúsculas. Cada forma em português tem
um sinônimo em inglês, e as duas podem ser usadas.

| Português | Inglês | | Português | Inglês |
| --- | --- | --- | --- | --- |
| `programa` | `program` | | `escreva` | `print` |
| `var` | `var` | | `leia` | `read` |
| `inteiro` | `int` | | `verdadeiro` | `true` |
| `real` | `float` | | `falso` | `false` |
| `booleano` | `boolean` | | `e` | `and` |
| `se` | `if` | | `ou` | `or` |
| `senão` | `else` | | `não` | `not` |
| `enquanto` | `while` | | `fim` | `end` |

## Identificadores

Começam com letra ou `_` e continuam com letras, dígitos ou `_`. As letras
aceitas são `a`–`z`, `A`–`Z`, `ã` e `Ã` (a última existe para reconhecer
`senão` e `não`). Outros acentos, como `ç` ou `é`, geram erro léxico.

## Exemplo completo

[`examples/programa.min`](../../minilang/examples/programa.min) lê um limite,
mostra os números pares até ele, a soma e a média:

```text
programa demonstracao {
    var limite: inteiro;
    var i: inteiro;
    var soma: inteiro;
    var media: real;
    var quantidade: inteiro;

    leia(limite);
    i = 1;
    enquanto (i <= limite) {
        se (i % 2 == 0) {
            escreva(i);
            quantidade = quantidade + 1;
        }
        soma = soma + i;
        i = i + 1;
    }

    escreva(soma);
    se (quantidade > 0) {
        media = soma;
        media = media / quantidade;
        escreva(media);
    } senão escreva(0);
} fim
```

Com a entrada `6`, imprime `2`, `4`, `6`, `21` e `7.0`.

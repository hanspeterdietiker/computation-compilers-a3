# Analisador semântico (M3)

O analisador semântico recebe a AST produzida pelo parser e percorre os comandos
em ordem, mantendo uma tabela de símbolos. A gramática livre de contexto
reconhece a forma dos comandos; regras como “identificador precisa estar
declarado”, “atribuição precisa ser compatível” e “condição precisa ser
booleana” dependem do contexto do programa. O ambiente de símbolos e a
verificação de tipos implementam essas restrições sensíveis ao contexto.

## Tabela de símbolos e escopos

Cada símbolo armazena nome, tipo, escopo e token da declaração (linha e
coluna). O escopo global contém as declarações do programa; cada bloco `{ ... }`
cria um escopo filho. A resolução procura primeiro no escopo atual e depois nos
escopos ancestrais. Assim, blocos podem sombrear variáveis externas, mas uma
segunda declaração no mesmo escopo é rejeitada. O escopo do bloco é removido
ao terminar sua análise.

As declarações são visíveis a partir de seu ponto no programa: uma referência
anterior à declaração não é resolvida. Declarações com erro de redeclaração
não substituem o símbolo original.

## Regras de tipos

* Literais inteiros, reais e booleanos inferem `inteiro`, `real` e `booleano`.
* `+`, `-`, `*` e `/` aceitam operandos numéricos; o resultado é `real` se
  qualquer operando for `real`, e `inteiro` caso contrário.
* `%` exige dois operandos `inteiro`.
* `<`, `<=`, `>` e `>=` aceitam operandos numéricos e resultam em `booleano`.
* `==` e `!=` aceitam dois operandos numéricos ou dois booleanos.
* `e`, `ou` e `não` exigem operandos booleanos.
* Condições de `se` e `enquanto` precisam ser booleanas.
* A atribuição aceita tipos iguais e promove `inteiro` para `real`; não há
  conversão implícita de `real` para `inteiro`.
* `leia` exige que o identificador esteja declarado.

Cada nó de expressão da AST expõe `GetType()` com o tipo inferido; expressões
que contêm erros permanecem com tipo `unknown` para evitar diagnósticos em
cascata. A análise acumula erros com fase, linha e coluna e continua nos demais
comandos.

## Execução

Dentro de `minilang/`, execute:

```bash
go run ./cmd/minilang check examples/parser/valido.min
```

O modo `check` executa lexer, parser e analisador semântico, e em caso de
sucesso imprime a tabela com identificador, tipo, escopo e posição de
declaração. O comando `go test ./semantic` executa os testes da etapa.

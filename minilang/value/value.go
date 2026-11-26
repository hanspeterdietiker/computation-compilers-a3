// Package value implementa os valores da MiniLang (int64, float64, bool) e suas operações.
package value

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"minilang/ast"
	"minilang/lexer"
)

// ErrDivisionByZero indica divisor zero em / ou %.
var ErrDivisionByZero = errors.New("divisão por zero")

// Binary aplica um operador binário, promovendo inteiro para real se preciso.
func Binary(op lexer.TokenType, left, right interface{}) (interface{}, error) {
	switch op {
	case lexer.TokenAnd:
		return asBool(left) && asBool(right), nil
	case lexer.TokenOr:
		return asBool(left) || asBool(right), nil
	case lexer.TokenEqual, lexer.TokenNotEqual:
		equal := equals(left, right)
		if op == lexer.TokenNotEqual {
			return !equal, nil
		}
		return equal, nil
	}

	leftInt, leftIsInt := left.(int64)
	rightInt, rightIsInt := right.(int64)
	if leftIsInt && rightIsInt {
		return integerBinary(op, leftInt, rightInt)
	}
	return realBinary(op, toReal(left), toReal(right))
}

func integerBinary(op lexer.TokenType, left, right int64) (interface{}, error) {
	switch op {
	case lexer.TokenPlus:
		return left + right, nil
	case lexer.TokenMinus:
		return left - right, nil
	case lexer.TokenMultiply:
		return left * right, nil
	case lexer.TokenDivide, lexer.TokenPercent:
		if right == 0 {
			return nil, ErrDivisionByZero
		}
		if op == lexer.TokenDivide {
			return left / right, nil
		}
		return left % right, nil
	}
	return compare(op, left, right)
}

func realBinary(op lexer.TokenType, left, right float64) (interface{}, error) {
	switch op {
	case lexer.TokenPlus:
		return left + right, nil
	case lexer.TokenMinus:
		return left - right, nil
	case lexer.TokenMultiply:
		return left * right, nil
	case lexer.TokenDivide:
		if right == 0 {
			return nil, ErrDivisionByZero
		}
		return left / right, nil
	}
	return compare(op, left, right)
}

// compare aplica um operador relacional.
func compare[T int64 | float64](op lexer.TokenType, left, right T) (interface{}, error) {
	switch op {
	case lexer.TokenLessThan:
		return left < right, nil
	case lexer.TokenLessEqual:
		return left <= right, nil
	case lexer.TokenGreaterThan:
		return left > right, nil
	case lexer.TokenGreaterEqual:
		return left >= right, nil
	}
	return nil, fmt.Errorf("operador '%s' não suportado", op)
}

// Unary aplica não, - ou +.
func Unary(op lexer.TokenType, operand interface{}) (interface{}, error) {
	switch op {
	case lexer.TokenNot:
		return !asBool(operand), nil
	case lexer.TokenMinus:
		switch v := operand.(type) {
		case int64:
			return -v, nil
		case float64:
			return -v, nil
		}
	case lexer.TokenPlus:
		switch operand.(type) {
		case int64, float64:
			return operand, nil
		}
	}
	return nil, fmt.Errorf("operador unário '%s' não suportado para %v", op, operand)
}

// Coerce aplica a única conversão implícita: inteiro → real.
func Coerce(target ast.Type, v interface{}) interface{} {
	if i, ok := v.(int64); ok && target == ast.TypeReal {
		return float64(i)
	}
	return v
}

// Zero é o valor inicial de uma variável declarada.
func Zero(typ ast.Type) interface{} {
	switch typ {
	case ast.TypeReal:
		return 0.0
	case ast.TypeBoolean:
		return false
	default:
		return int64(0)
	}
}

// Format formata para escreva; reais usam 15 dígitos significativos e sempre têm parte decimal.
func Format(v interface{}) string {
	switch v := v.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		v, _ = strconv.ParseFloat(strconv.FormatFloat(v, 'g', 15, 64), 64)
		text := strconv.FormatFloat(v, 'f', -1, 64)
		if !math.IsInf(v, 0) && !math.IsNaN(v) && v == math.Trunc(v) {
			text += ".0"
		}
		return text
	case bool:
		if v {
			return "verdadeiro"
		}
		return "falso"
	}
	return fmt.Sprint(v)
}

// Parse converte a entrada de leia; real aceita forma inteira.
func Parse(typ ast.Type, text string) (interface{}, error) {
	invalid := fmt.Errorf("valor '%s' inválido para %s", text, typ)
	switch typ {
	case ast.TypeInteger:
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, invalid
		}
		return v, nil
	case ast.TypeReal:
		v, err := strconv.ParseFloat(text, 64)
		if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
			return nil, invalid
		}
		return v, nil
	case ast.TypeBoolean:
		switch text {
		case "verdadeiro", "true":
			return true, nil
		case "falso", "false":
			return false, nil
		}
	}
	return nil, invalid
}

func equals(left, right interface{}) bool {
	if l, ok := left.(bool); ok {
		r, ok := right.(bool)
		return ok && l == r
	}
	if l, ok := left.(int64); ok {
		if r, ok := right.(int64); ok {
			return l == r
		}
	}
	return toReal(left) == toReal(right)
}

func toReal(v interface{}) float64 {
	switch v := v.(type) {
	case int64:
		return float64(v)
	case float64:
		return v
	}
	return math.NaN()
}

func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

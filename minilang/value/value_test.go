package value

import (
	"errors"
	"testing"

	"minilang/ast"
	"minilang/lexer"
)

func TestBinaryPromotesIntegerToReal(t *testing.T) {
	cases := []struct {
		name        string
		op          lexer.TokenType
		left, right interface{}
		want        interface{}
	}{
		{"inteiro+inteiro", lexer.TokenPlus, int64(10), int64(5), int64(15)},
		{"inteiro+real", lexer.TokenPlus, int64(10), 2.5, 12.5},
		{"real*inteiro", lexer.TokenMultiply, 19.9, int64(2), 39.8},
		{"divisão inteira trunca", lexer.TokenDivide, int64(7), int64(2), int64(3)},
		{"divisão real", lexer.TokenDivide, 7.0, int64(2), 3.5},
		{"módulo", lexer.TokenPercent, int64(7), int64(3), int64(1)},
		{"menor misto", lexer.TokenLessThan, int64(2), 2.5, true},
		{"igualdade mista", lexer.TokenEqual, int64(2), 2.0, true},
		{"diferença booleana", lexer.TokenNotEqual, true, false, true},
		{"e lógico", lexer.TokenAnd, true, false, false},
		{"ou lógico", lexer.TokenOr, false, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Binary(tc.op, tc.left, tc.right)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if got != tc.want {
				t.Fatalf("Binary(%s, %v, %v) = %#v; esperado %#v", tc.op, tc.left, tc.right, got, tc.want)
			}
		})
	}
}

func TestBinaryRejectsDivisionByZero(t *testing.T) {
	for _, op := range []lexer.TokenType{lexer.TokenDivide, lexer.TokenPercent} {
		if _, err := Binary(op, int64(1), int64(0)); !errors.Is(err, ErrDivisionByZero) {
			t.Errorf("%s por zero inteiro: erro = %v", op, err)
		}
	}
	if _, err := Binary(lexer.TokenDivide, 1.0, int64(0)); !errors.Is(err, ErrDivisionByZero) {
		t.Errorf("divisão real por zero: erro = %v", err)
	}
}

func TestUnary(t *testing.T) {
	if got, _ := Unary(lexer.TokenMinus, int64(3)); got != int64(-3) {
		t.Errorf("-3 = %#v", got)
	}
	if got, _ := Unary(lexer.TokenMinus, 2.5); got != -2.5 {
		t.Errorf("-2.5 = %#v", got)
	}
	if got, _ := Unary(lexer.TokenNot, true); got != false {
		t.Errorf("não verdadeiro = %#v", got)
	}
}

func TestCoerceAndZero(t *testing.T) {
	if got := Coerce(ast.TypeReal, int64(10)); got != 10.0 {
		t.Errorf("Coerce(real, 10) = %#v", got)
	}
	if got := Zero(ast.TypeReal); got != 0.0 {
		t.Errorf("Zero(real) = %#v", got)
	}
	if got := Zero(ast.TypeBoolean); got != false {
		t.Errorf("Zero(booleano) = %#v", got)
	}
}

func TestFormat(t *testing.T) {
	cases := map[interface{}]string{
		int64(42): "42",
		12.5:      "12.5",
		10.0:      "10.0",
		-0.25:     "-0.25",
		true:      "verdadeiro",
		false:     "falso",
	}
	for input, want := range cases {
		if got := Format(input); got != want {
			t.Errorf("Format(%#v) = %q; esperado %q", input, got, want)
		}
	}
}

func TestParse(t *testing.T) {
	if got, err := Parse(ast.TypeInteger, "15"); err != nil || got != int64(15) {
		t.Errorf("Parse(inteiro, 15) = %#v, %v", got, err)
	}
	if got, err := Parse(ast.TypeReal, "3"); err != nil || got != 3.0 {
		t.Errorf("Parse(real, 3) = %#v, %v", got, err)
	}
	if got, err := Parse(ast.TypeBoolean, "verdadeiro"); err != nil || got != true {
		t.Errorf("Parse(booleano, verdadeiro) = %#v, %v", got, err)
	}
	if _, err := Parse(ast.TypeInteger, "2.5"); err == nil {
		t.Error("Parse(inteiro, 2.5) deveria falhar")
	}
}

func TestFormatHidesFloatingPointNoise(t *testing.T) {
	product, _ := Binary(lexer.TokenMultiply, 19.9, int64(3))
	if got := Format(product); got != "59.7" {
		t.Errorf("Format(19.9 * 3) = %q; esperado \"59.7\"", got)
	}
}

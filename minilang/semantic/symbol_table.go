package semantic

import (
	"fmt"

	"minilang/ast"
	"minilang/lexer"
)

type Scope struct {
	ID       int
	Name     string
	ParentID int
	Depth    int
}

type Symbol struct {
	Name             string
	Type             ast.Type
	Scope            string
	ScopeID          int
	DeclarationToken lexer.Token
}

type SymbolTable struct {
	scopes      []Scope
	stack       []int
	symbols     []Symbol
	byScope     []map[string]Symbol
	nextScopeID int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		scopes:  []Scope{{ID: 0, Name: "global", Depth: 0}},
		stack:   []int{0},
		byScope: []map[string]Symbol{{}},
	}
}

func (t *SymbolTable) CurrentScope() Scope {
	return t.scopes[t.stack[len(t.stack)-1]]
}

func (t *SymbolTable) PushScope(token lexer.Token) Scope {
	parent := t.CurrentScope()
	t.nextScopeID++
	scope := Scope{
		ID:       t.nextScopeID,
		Name:     fmt.Sprintf("bloco@%d:%d", token.Line, token.Column),
		ParentID: parent.ID,
		Depth:    parent.Depth + 1,
	}
	t.scopes = append(t.scopes, scope)
	t.byScope = append(t.byScope, make(map[string]Symbol))
	t.stack = append(t.stack, scope.ID)
	return scope
}

func (t *SymbolTable) PopScope() {
	if len(t.stack) > 1 {
		t.stack = t.stack[:len(t.stack)-1]
	}
}

func (t *SymbolTable) Declare(name string, typ ast.Type, token lexer.Token) (Symbol, bool) {
	scope := t.CurrentScope()
	if _, exists := t.byScope[scope.ID][name]; exists {
		return Symbol{}, false
	}
	symbol := Symbol{
		Name:             name,
		Type:             typ,
		Scope:            scope.Name,
		ScopeID:          scope.ID,
		DeclarationToken: token,
	}
	t.byScope[scope.ID][name] = symbol
	t.symbols = append(t.symbols, symbol)
	return symbol, true
}

func (t *SymbolTable) Resolve(name string) (Symbol, bool) {
	for i := len(t.stack) - 1; i >= 0; i-- {
		if symbol, exists := t.byScope[t.stack[i]][name]; exists {
			return symbol, true
		}
	}
	return Symbol{}, false
}

func (t *SymbolTable) Symbols() []Symbol {
	return append([]Symbol(nil), t.symbols...)
}

func (t *SymbolTable) Scopes() []Scope {
	return append([]Scope(nil), t.scopes...)
}

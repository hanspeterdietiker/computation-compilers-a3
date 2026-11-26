package interpreter

import "minilang/ast"

// variable guarda o tipo declarado para coerção e leitura.
type variable struct {
	typ   ast.Type
	value interface{}
}

// environment é o escopo de execução; cada bloco cria um filho.
type environment struct {
	parent    *environment
	variables map[string]*variable
}

func newEnvironment(parent *environment) *environment {
	return &environment{parent: parent, variables: make(map[string]*variable)}
}

func (e *environment) declare(name string, typ ast.Type, initial interface{}) {
	e.variables[name] = &variable{typ: typ, value: initial}
}

func (e *environment) resolve(name string) (*variable, bool) {
	for env := e; env != nil; env = env.parent {
		if found, ok := env.variables[name]; ok {
			return found, true
		}
	}
	return nil, false
}

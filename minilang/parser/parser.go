package parser

import (
	"fmt"
	"strconv"

	"minilang/ast"
	"minilang/lexer"
)

type ParseError struct {
	Token   lexer.Token
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[ERRO SINTÁTICO] linha %d, coluna %d: %s", e.Token.Line, e.Token.Column, e.Message)
}

type Parser struct {
	tokens  []lexer.Token
	current int
	errors  []error
}

func New(tokens []lexer.Token) *Parser {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != lexer.TokenEOF {
		tokens = append(append([]lexer.Token(nil), tokens...), lexer.Token{Type: lexer.TokenEOF})
	}
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*ast.Program, []error) {
	p.current = 0
	p.errors = nil
	programToken := p.peek()
	program := &ast.Program{NodeBase: ast.NodeBase{Token: programToken}}

	if p.match(lexer.TokenProgram) {
		programToken = p.previous()
		program.Token = programToken
		if p.match(lexer.TokenIdentifier) {
			name := p.previous()
			program.Name = &name
		}
	}

	hasBraces := p.match(lexer.TokenLeftBrace)
	for !p.check(lexer.TokenEOF) &&
		!(hasBraces && p.check(lexer.TokenRightBrace)) &&
		!(!hasBraces && p.check(lexer.TokenEnd)) {
		before := p.current
		if statement := p.parseStatementSafely(); statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		if p.current == before && !p.check(lexer.TokenEOF) &&
			!(hasBraces && p.check(lexer.TokenRightBrace)) {
			p.advance()
		}
	}

	if hasBraces && !p.match(lexer.TokenRightBrace) {
		p.errors = append(p.errors, ParseError{Token: p.peek(), Message: "esperado '}' para fechar o programa"})
	}
	if p.match(lexer.TokenEnd) {
		if program.Token.Type == lexer.TokenEOF {
			program.Token = p.previous()
		}
	}
	if !p.check(lexer.TokenEOF) {
		p.errors = append(p.errors, ParseError{Token: p.peek(), Message: "token inesperado após o fim do programa"})
	}
	return program, p.errors
}

func (p *Parser) parseStatementSafely() (statement ast.Statement) {
	defer func() {
		if recovered := recover(); recovered != nil {
			parseError, ok := recovered.(ParseError)
			if !ok {
				panic(recovered)
			}
			p.errors = append(p.errors, parseError)
			p.synchronize()
		}
	}()
	return p.parseStatement()
}

func (p *Parser) parseStatement() ast.Statement {
	if p.match(lexer.TokenVar) {
		return p.parseVarDecl(p.previous())
	}
	if p.match(lexer.TokenIf) {
		return p.parseIf(p.previous())
	}
	if p.match(lexer.TokenWhile) {
		return p.parseWhile(p.previous())
	}
	if p.match(lexer.TokenPrint) {
		return p.parsePrint(p.previous())
	}
	if p.match(lexer.TokenRead) {
		return p.parseRead(p.previous())
	}
	if p.match(lexer.TokenLeftBrace) {
		return p.parseBlock(p.previous())
	}
	if p.check(lexer.TokenIdentifier) {
		name := p.advance()
		p.consume(lexer.TokenAssign, "esperado '=' após o identificador")
		value := p.parseExpression()
		p.consume(lexer.TokenSemicolon, "esperado ';' após atribuição")
		return &ast.AssignStmt{NodeBase: ast.NodeBase{Token: name}, Name: name, Value: value}
	}
	p.fail(p.peek(), "comando esperado")
	return nil
}

func (p *Parser) parseVarDecl(start lexer.Token) ast.Statement {
	name := p.consume(lexer.TokenIdentifier, "esperado identificador após 'var'")
	p.consume(lexer.TokenColon, "esperado ':' após o nome da variável")
	typeToken := p.peek()
	if !p.match(lexer.TokenTypeInt, lexer.TokenTypeFloat, lexer.TokenTypeBool) {
		p.fail(p.peek(), "esperado tipo 'int', 'float' ou 'boolean'")
	}
	p.consume(lexer.TokenSemicolon, "esperado ';' após declaração")
	return &ast.VarDecl{NodeBase: ast.NodeBase{Token: start}, Name: name, Type: typeToken}
}

func (p *Parser) parseIf(start lexer.Token) ast.Statement {
	p.consume(lexer.TokenLeftParen, "esperado '(' após 'if'")
	condition := p.parseExpression()
	p.consume(lexer.TokenRightParen, "esperado ')' após condição")
	thenBranch := p.parseStatement()
	var elseBranch ast.Statement
	if p.match(lexer.TokenElse) {
		elseBranch = p.parseStatement()
	}
	return &ast.IfStmt{
		NodeBase: ast.NodeBase{Token: start}, Condition: condition,
		ThenBranch: thenBranch, ElseBranch: elseBranch,
	}
}

func (p *Parser) parseWhile(start lexer.Token) ast.Statement {
	p.consume(lexer.TokenLeftParen, "esperado '(' após 'while'")
	condition := p.parseExpression()
	p.consume(lexer.TokenRightParen, "esperado ')' após condição")
	body := p.parseStatement()
	return &ast.WhileStmt{NodeBase: ast.NodeBase{Token: start}, Condition: condition, Body: body}
}

func (p *Parser) parsePrint(start lexer.Token) ast.Statement {
	p.consume(lexer.TokenLeftParen, "esperado '(' após 'print'")
	expression := p.parseExpression()
	p.consume(lexer.TokenRightParen, "esperado ')' após expressão")
	p.consume(lexer.TokenSemicolon, "esperado ';' após 'print'")
	return &ast.PrintStmt{NodeBase: ast.NodeBase{Token: start}, Expression: expression}
}

func (p *Parser) parseRead(start lexer.Token) ast.Statement {
	p.consume(lexer.TokenLeftParen, "esperado '(' após 'read'")
	name := p.consume(lexer.TokenIdentifier, "esperado identificador em 'read'")
	p.consume(lexer.TokenRightParen, "esperado ')' após identificador")
	p.consume(lexer.TokenSemicolon, "esperado ';' após 'read'")
	return &ast.ReadStmt{NodeBase: ast.NodeBase{Token: start}, Name: name}
}

func (p *Parser) parseBlock(start lexer.Token) ast.Statement {
	block := &ast.BlockStmt{NodeBase: ast.NodeBase{Token: start}}
	for !p.check(lexer.TokenRightBrace) && !p.check(lexer.TokenEOF) {
		before := p.current
		if statement := p.parseStatementSafely(); statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		if p.current == before && !p.check(lexer.TokenRightBrace) && !p.check(lexer.TokenEOF) {
			p.advance()
		}
	}
	p.consume(lexer.TokenRightBrace, "esperado '}' para fechar bloco")
	return block
}

func (p *Parser) parseExpression() ast.Expression { return p.parseOr() }

func (p *Parser) parseOr() ast.Expression {
	expression := p.parseAnd()
	for p.match(lexer.TokenOr) {
		operator := p.previous()
		right := p.parseAnd()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseAnd() ast.Expression {
	expression := p.parseEquality()
	for p.match(lexer.TokenAnd) {
		operator := p.previous()
		right := p.parseEquality()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseEquality() ast.Expression {
	expression := p.parseComparison()
	for p.match(lexer.TokenEqual, lexer.TokenNotEqual) {
		operator := p.previous()
		right := p.parseComparison()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseComparison() ast.Expression {
	expression := p.parseTerm()
	for p.match(lexer.TokenLessThan, lexer.TokenLessEqual, lexer.TokenGreaterThan, lexer.TokenGreaterEqual) {
		operator := p.previous()
		right := p.parseTerm()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseTerm() ast.Expression {
	expression := p.parseFactor()
	for p.match(lexer.TokenPlus, lexer.TokenMinus) {
		operator := p.previous()
		right := p.parseFactor()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseFactor() ast.Expression {
	expression := p.parseUnary()
	for p.match(lexer.TokenMultiply, lexer.TokenDivide, lexer.TokenPercent) {
		operator := p.previous()
		right := p.parseUnary()
		expression = &ast.BinaryExpr{NodeBase: ast.NodeBase{Token: operator}, Left: expression, Operator: operator, Right: right}
	}
	return expression
}

func (p *Parser) parseUnary() ast.Expression {
	if p.match(lexer.TokenNot, lexer.TokenMinus, lexer.TokenPlus) {
		operator := p.previous()
		return &ast.UnaryExpr{NodeBase: ast.NodeBase{Token: operator}, Operator: operator, Right: p.parseUnary()}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expression {
	if p.match(lexer.TokenFalse) {
		return &ast.LiteralExpr{NodeBase: ast.NodeBase{Token: p.previous()}, Value: false}
	}
	if p.match(lexer.TokenTrue) {
		return &ast.LiteralExpr{NodeBase: ast.NodeBase{Token: p.previous()}, Value: true}
	}
	if p.match(lexer.TokenIntLiteral) {
		token := p.previous()
		value, err := strconv.ParseInt(token.Lexeme, 10, 64)
		if err != nil {
			p.fail(token, "literal inteiro fora do intervalo permitido")
		}
		return &ast.LiteralExpr{NodeBase: ast.NodeBase{Token: token}, Value: value}
	}
	if p.match(lexer.TokenFloatLiteral) {
		token := p.previous()
		value, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			p.fail(token, "literal real inválido")
		}
		return &ast.LiteralExpr{NodeBase: ast.NodeBase{Token: token}, Value: value}
	}
	if p.match(lexer.TokenIdentifier) {
		name := p.previous()
		return &ast.IdentifierExpr{NodeBase: ast.NodeBase{Token: name}, Name: name}
	}
	if p.match(lexer.TokenLeftParen) {
		start := p.previous()
		expression := p.parseExpression()
		p.consume(lexer.TokenRightParen, "esperado ')' após expressão")
		return &ast.GroupingExpr{NodeBase: ast.NodeBase{Token: start}, Expression: expression}
	}
	p.fail(p.peek(), "expressão esperada")
	return nil
}

func (p *Parser) synchronize() {
	for !p.check(lexer.TokenEOF) {
		if p.previous().Type == lexer.TokenSemicolon || p.check(lexer.TokenRightBrace) {
			return
		}
		switch p.peek().Type {
		case lexer.TokenVar, lexer.TokenIf, lexer.TokenWhile, lexer.TokenPrint,
			lexer.TokenRead, lexer.TokenLeftBrace, lexer.TokenIdentifier:
			return
		}
		p.advance()
	}
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.fail(p.peek(), message)
	return lexer.Token{}
}

func (p *Parser) fail(token lexer.Token, message string) {
	panic(ParseError{Token: token, Message: message})
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType lexer.TokenType) bool { return p.peek().Type == tokenType }

func (p *Parser) advance() lexer.Token {
	if !p.check(lexer.TokenEOF) {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool { return p.check(lexer.TokenEOF) }

func (p *Parser) peek() lexer.Token { return p.tokens[p.current] }

func (p *Parser) previous() lexer.Token {
	if p.current == 0 {
		return p.tokens[0]
	}
	return p.tokens[p.current-1]
}

package lexer

import (
	"fmt"
)

type Lexer struct {
	source  []rune
	start   int
	current int
	line    int
	column  int
	tokens  []Token
}

func New(source string) *Lexer {
	return &Lexer{
		source: []rune(source),
		line:   1,
		column: 1,
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}
func (l *Lexer) advance() rune {
	c := l.source[l.current]
	l.current++
	l.column++

	return c
}

func (l *Lexer) scanToken() {
	c := l.advance()

	switch c {

	case '+':
		l.addToken(TokenPlus)

	case '-':
		l.addToken(TokenMinus)

	case '*':
		l.addToken(TokenMultiply)

	case '/':
		l.addToken(TokenDivide)

	case '#':
		l.addToken(TokenComment)
		l.comment()

	case '%':
		l.addToken(TokenPercent)

	case '(':
		l.addToken(TokenLeftParen)

	case ')':
		l.addToken(TokenRightParen)

	case '{':
		l.addToken(TokenLeftBrace)

	case '}':
		l.addToken(TokenRightBrace)

	case ';':
		l.addToken(TokenSemicolon)

	case ':':
		l.addToken(TokenColon)

	case ',':
		l.addToken(TokenComma)

	case '.':
		l.addToken(TokenDot)

	case ' ', '\r', '\t':
		// ignora espaço em branco

	case '\n':
		l.line++
		l.column = 1

	case '=':
		if l.match('=') {
			l.addToken(TokenEqual)
		} else {
			l.addToken(TokenAssign)
		}
	case '<':
		if l.match('=') {
			l.addToken(TokenLessEqual)
		} else {
			l.addToken(TokenLessThan)
		}

	case '>':
		if l.match('=') {
			l.addToken(TokenGreaterEqual)
		} else {
			l.addToken(TokenGreaterThan)
		}

	case '!':
		if l.match('=') {
			l.addToken(TokenNotEqual)
		} else {
			l.reportInvalidCharacter(c)
		}
	default:
		if isDigit(c) {
			l.number()
		} else if isLetter(c) {
			l.identifier()
		} else {
			l.reportInvalidCharacter(c)
		}

	}

}

func (l *Lexer) number() {

	for !l.isAtEnd() && isDigit(l.source[l.current]) {
		l.advance()
	}

	if !l.isAtEnd() &&
		l.source[l.current] == '.' &&
		l.peekNextIsDigit() {

		l.advance()

		for !l.isAtEnd() && isDigit(l.source[l.current]) {
			l.advance()
		}

		l.addToken(TokenFloatLiteral)
		return
	}

	l.addToken(TokenIntLiteral)
}

func (l *Lexer) comment() {
	for !l.isAtEnd() && l.source[l.current] != '\n' {
		l.advance()
	}
}

func (l *Lexer) reportInvalidCharacter(c rune) {
	fmt.Printf(
		"[ERRO LÉXICO] linha %d, coluna %d: caractere '%c' não reconhecido\n",
		l.line,
		l.column-1,
		c,
	)
}

func (l *Lexer) identifier() {

	for !l.isAtEnd() && isAlphaNumeric(l.source[l.current]) {
		l.advance()
	}

	text := string(l.source[l.start:l.current])

	tokenType, exists := Keywords[text]

	if !exists {
		tokenType = TokenIdentifier
	}

	l.addToken(tokenType)
}
func (l *Lexer) peekNextIsDigit() bool {

	if l.current+1 >= len(l.source) {
		return false
	}

	return isDigit(l.source[l.current+1])
}
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_' ||
		c == 'ã' ||
		c == 'Ã'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c rune) bool {
	return isLetter(c) || isDigit(c)
}

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() {
		return false
	}

	if l.source[l.current] != expected {
		return false
	}

	l.current++
	l.column++

	return true
}

func (l *Lexer) addToken(tokenType TokenType) {
	text := string(l.source[l.start:l.current])

	l.tokens = append(l.tokens, Token{
		Type:   tokenType,
		Lexeme: text,
		Line:   l.line,
		Column: l.column - len([]rune(text)),
	})
}

func (l *Lexer) ScanTokens() []Token {

	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.tokens = append(l.tokens, Token{
		Type:   TokenEOF,
		Lexeme: "",
		Line:   l.line,
		Column: l.column,
	})
	return l.tokens
}

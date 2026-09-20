package lexer

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
	}
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

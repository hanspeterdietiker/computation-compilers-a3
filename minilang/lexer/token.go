package lexer

type TokenType string

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

const (
	// Literals and Identifiers
	TokenIdentifier   TokenType = "IDENTIFIER"
	TokenIntLiteral   TokenType = "INT_LITERAL"
	TokenFloatLiteral TokenType = "FLOAT_LITERAL"

	// Types
	TokenTypeInt   TokenType = "TYPE_INT"
	TokenTypeFloat TokenType = "TYPE_FLOAT"
	TokenTypeBool  TokenType = "TYPE_BOOL"

	// Conditionals
	TokenProgram TokenType = "PROGRAM"
	TokenVar     TokenType = "VAR"

	TokenIf    TokenType = "IF"
	TokenElse  TokenType = "ELSE"
	TokenWhile TokenType = "WHILE"

	TokenPrint TokenType = "PRINT"
	TokenRead  TokenType = "READ"

	TokenTrue  TokenType = "TRUE"
	TokenFalse TokenType = "FALSE"

	TokenAnd TokenType = "AND"
	TokenOr  TokenType = "OR"
	TokenNot TokenType = "NOT"

	TokenEnd TokenType = "END"

	// Arithmetic Operators
	TokenPlus     TokenType = "PLUS"     // -> +
	TokenMinus    TokenType = "MINUS"    // -> -
	TokenMultiply TokenType = "MULTIPLY" // -> *
	TokenDivide   TokenType = "DIVIDE"   // -> /
	TokenPercent  TokenType = "PERCENT"  // -> %

	// Relational Operators
	TokenAssign       TokenType = "ASSIGN"        // -> =
	TokenEqual        TokenType = "EQUAL"         // -> ==
	TokenNotEqual     TokenType = "NOT_EQUAL"     // -> !=
	TokenLessThan     TokenType = "LESS_THAN"     // -> <
	TokenGreaterThan  TokenType = "GREATER_THAN"  // -> >
	TokenLessEqual    TokenType = "LESS_EQUAL"    // -> <=
	TokenGreaterEqual TokenType = "GREATER_EQUAL" // -> >=

	// Delimiters
	TokenLeftParen  TokenType = "LEFT_PAREN"  // -> (
	TokenRightParen TokenType = "RIGHT_PAREN" // -> )

	TokenLeftBrace  TokenType = "LEFT_BRACE"  // -> {
	TokenRightBrace TokenType = "RIGHT_BRACE" // -> }

	TokenSemicolon TokenType = "SEMICOLON" // -> ;
	TokenComma     TokenType = "COMMA"     // -> ,
	TokenDot       TokenType = "DOT"       // -> .
	TokenColon     TokenType = "COLON"     // -> :

	TokenComment TokenType = "COMMENT" // -> #


	// Eof
	TokenEOF TokenType = "EOF" // End of file
)

var Keywords = map[string]TokenType{
	"program": TokenProgram,
	"var":     TokenVar,

	"int":     TokenTypeInt,
	"float":   TokenTypeFloat,
	"boolean": TokenTypeBool,

	"if":    TokenIf,
	"else":  TokenElse,
	"while": TokenWhile,

	"print": TokenPrint,
	"read":  TokenRead,

	"true":  TokenTrue,
	"false": TokenFalse,

	"and": TokenAnd,
	"or":  TokenOr,
	"not": TokenNot,

	"end": TokenEnd,
}

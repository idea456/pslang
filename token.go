package main

type TokenType int

const (
	IDENTIFIER TokenType = iota
	STRING
	NUMBER

	AND
	CLASS
	ELSE
	FALSE
	PROCEDURE
	FOR
	IF
	THEN
	NIL
	OR
	SAY
	RETURN
	PARENT
	THIS
	TRUE
	WHILE
	DO
	ASSUME
	SET
	TO

	NOT           // !
	NOT_EQUAL     // !=
	EQUAL_EQUAL   // ==
	LESS          // <
	LESS_EQUAL    // <=
	GREATER       // >
	GREATER_EQUAL // >=

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	MODULUS
	INCREMENT
	DECREMENT
	BY

	EMPTY
	UNKNOWN
	EOF
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

var keywords = map[string]TokenType{
	"and":       AND,
	"class":     CLASS,
	"else":      ELSE,
	"false":     FALSE,
	"procedure": PROCEDURE,
	"for":       FOR,
	"if":        IF,
	"then":      THEN,
	"or":        OR,
	"say":       SAY,
	"return":    RETURN,
	"parent":    PARENT, // TODO: can change to PARENT
	"this":      THIS,
	"true":      TRUE,
	"while":     WHILE,
	"do":        DO,
	"assume":    ASSUME,
	"set":       SET,
	"to":        TO,
	"increment": INCREMENT,
	"decrement": DECREMENT,
	"by":        BY,
	"!":         NOT,
	"!=":        NOT_EQUAL,
	"==":        EQUAL_EQUAL,
	"<":         LESS,
	"<=":        LESS_EQUAL,
	">":         GREATER,
	">=":        GREATER_EQUAL,
	"(":         LEFT_PAREN,
	")":         RIGHT_PAREN,
	"{":         LEFT_BRACE,
	"}":         RIGHT_BRACE,
	",":         COMMA,
	".":         DOT,
	"-":         MINUS,
	"+":         PLUS,
	";":         SEMICOLON,
	"/":         SLASH,
	"*":         STAR,
	"%":         MODULUS,
	"empty":     EMPTY,
}

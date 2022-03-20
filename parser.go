package main

import "fmt"

/*
Stratified grammar:

file -> declaration* EOF;
declaration -> var_declaration | statement;
var_declaration -> "set" IDENTIFIER ("to" (expression | incr_decr))? ";"
statement -> say_stmt | expr_stmt | incr_decr_stmt | if_stmt | while_stmt | block_stmt;
say_stmt -> "say" expression ";"
expr_stmt -> expression ";"
incr_decr_stmt -> ("increment" | "decrement") IDENTIFIER "by" expression;
if_stmt -> "if" expression "then" statement ("else" statement)?;
while_stmt -> "while" expression "do" "{" statement "}";
block_stmt -> "{" declaration* "}"

expression -> equality | logical_or;
logical_or -> logical_and ("or" logical_and)*;
logical_and -> equality ("and" equality)*;
equality -> comparison (("==" | "!=") comparison)*;
comparison -> term ((">" | ">=" | "<" | "<=") term)*;
term -> factor (("+" | "-") factor)*;
factor -> unary (("*" | "/" | "") unary)*
unary -> ("-" | "!") unary | primary;
primary -> NUMBER | STRING | IDENTIFIER | "true" | "false" | "empty" | "(" expression ")";
*/

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	var parser Parser = Parser{}
	parser.tokens = tokens
	parser.current = 0
	return &parser
}

func (p *Parser) Parse() []Statement {
	defer func() {
		// exit panic mode and synchronize to the nearest starting statement keyword
		if r := recover(); r != nil {
			fmt.Printf("%+v\n", r)
			p.synchronize()
			// if p.peek().tokenType == LEFT_BRACE {
			// 	p.synchronizeToken(RIGHT_BRACE)
			// 	fmt.Println(p.peek().lexeme)
			// 	p.next()
			// }
			p.Parse()
		}
	}()
	var statements []Statement = make([]Statement, 0)
	// var expr Expression = p.expression()
	

	for p.peek().tokenType != EOF {
		statements = append(statements, p.declaration())
	}

	return statements
}

/*
declaration -> var_declaration | statement;
*/
func (p *Parser) declaration() Statement {
	// every statement must have terminating semicolon
	// defer func() {
	// 	if !p.end() && !p.match(SEMICOLON) {
	// 		// if p.previous().tokenType == SEMICOLON && p.peek().tokenType == RIGHT_BRACE {
	// 		// 	return
	// 		// }
	// 		if p.previous().tokenType == SEMICOLON && p.match(RIGHT_BRACE) {
	// 			return
	// 		}
	// 		RuntimeError(p.previous().line, p.previous().lexeme, "expected semicolon after statement!")
	// 	}
	// }()

	defer func() {
		// skip the NEWLINE token to go to next line (or else we end up in infinite loop if current token is NEWLINE)
		if !p.end() && p.match(NEWLINE) {
			return
		}
	}()

	if p.match(SET) {
		return p.var_declaration()
	}

	return p.statement()
}

/*
var_declaration -> "set" IDENTIFIER ("to" expression)? ";"
*/
func (p *Parser) var_declaration() Statement {
	var identifier Token = p.next()

	var expr Expression = nil
	for p.match(TO) {
		expr = p.expression()
	}

	return &VariableStmt{name: identifier, initializer: expr}
}

/*
statement -> say_stmt | expr_stmt | incr_decr_stmt | if_stmt | block;
*/
func (p *Parser) statement() Statement {
	if p.match(SAY) {
		return p.say_stmt()
	}
	if p.match(LEFT_BRACE) {
		return p.block_stmt()
	}
	if p.match(INCREMENT, DECREMENT) {
		return p.incr_decr_stmt()
	}
	if p.match(IF) {
		return p.if_stmt()
	}
	if p.match(WHILE) {
		return p.while_stmt()
	}
	return p.expr_stmt()

}

/*
say_stmt -> "say" expression ";"
*/
func (p *Parser) say_stmt() Statement {
	var expr Expression = p.expression()
	// expression returns EOF
	if expr == nil {
		RuntimeError(p.previous().line, p.previous().lexeme, "Expected expression after 'say' statement.")
	}

	return &SayStmt{
		expression: expr,
	}
}

/*
expr_stmt -> expression ";"
*/
func (p *Parser) expr_stmt() Statement {
	return &ExprStmt{
		expression: p.expression(), 
	}
}

/*
incr_decr_stmt -> ("increment" | "decrement") INDENTIFIER "by" expression;
*/
func (p *Parser) incr_decr_stmt() Statement {
	var operator Token = p.previous()
	var identifier Token = p.peek()

	p.next()

	if p.match(BY) {
		return &IncrDecrStmt{
			identifier: identifier,
			operator:   operator,
			right:      p.expression(),
		}
	} else {
		panic("Increment/decrement statements must be followed with 'by'.")
	}
}

func (p *Parser) if_stmt() Statement {
	// p.consume(LEFT_PAREN, "Error, expected '(' in if statement")
	var expr Expression = p.expression()
	// p.consume(RIGHT_PAREN, "Error, expected ')' after if statement")
	p.consume(THEN, "Error, if statements are followed by 'then'")
	p.consume(LEFT_BRACE, "Error, if statements are enclosed by braces'")

	var thenBranch Statement = p.statement()
	p.consume(RIGHT_BRACE, "Error, if statements are enclosed by braces")
	var elseBranch Statement
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{
		expression: expr,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

/*
while_stmt -> "while" expression "do" "{" statement "}";
*/
func (p *Parser) while_stmt() Statement {
	var expr Expression = p.expression()

	if !p.match(DO) {
		var token Token = p.peek()
		RuntimeError(token.line, token.lexeme, "expected 'do' after while statement.")
	}

	// if !p.match(LEFT_BRACE) {
	// 	var token Token = p.peek()
	// 	RuntimeError(token.line, token.lexeme, "expected '{' in while statement.")
	// }

	// for !p.match(RIGHT_BRACE) {
	// 	body := &
	// }
	// body := p.statement()

	// if p.end() || !p.match(RIGHT_BRACE) {
	// 	var token Token = p.peek()
	// 	RuntimeError(token.line, token.lexeme, "expected '}' in while statement.")
	// }

	return &WhileStmt{
		condition: expr,
		body:      p.statement(),
	}
}

/*
block -> "{" declaration* "}"
*/
func (p *Parser) block_stmt() Statement {
	var statements []Statement = make([]Statement, 0)

	for !(p.match(RIGHT_BRACE)) && p.peek().tokenType != EOF {
		statements = append(statements, p.declaration())
	}

	// if p.peek().tokenType == EOF && p.previous().tokenType != RIGHT_BRACE {
	// 	var token Token = p.peek()
	// 	RuntimeError(token.line, token.lexeme, "expect closing braces in block statement!")
	// }
	// !p.match(RIGHT_BRACE)
	if p.end() || p.previous().tokenType != RIGHT_BRACE {
		var token Token = p.peek()
		RuntimeError(token.line, token.lexeme, "expect closing braces in block statement!")
	}
	return &BlockStmt{
		statements: statements,
	}
}

func (p *Parser) expression() Expression {
	// if p.match(INCREMENT, DECREMENT) {
	// 	if p.match(BY) {

	// 	} else {
	// 		panic("Increment/decrement must follow 'by' keyword!")
	// 	}
	// }

	// return p.equality()
	return p.logical_or()
}

/*
logical_or -> logical_and ("or" logical_and)*;
logical_and -> equality ("and" equality)*;
*/
func (p *Parser) logical_or() Expression {
	var expr Expression = p.logical_and()

	for p.match(OR) {
		expr = &Logical{
			left:     expr,
			operator: p.previous(),
			right:    p.logical_and(),
		}
	}
	return expr
}

func (p *Parser) logical_and() Expression {
	var expr Expression = p.equality()

	for p.match(AND) {
		expr = &Logical{
			left:     expr,
			operator: p.previous(),
			right:    p.equality(),
		}
	}
	return expr
}

func (p *Parser) equality() Expression {
	var expr Expression = p.comparison()

	for p.match(EQUAL_EQUAL, NOT_EQUAL) {
		var operator Token = p.previous()
		var right Expression = p.comparison()
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

/*
comparison -> term ((">" | ">=" | "<" | "<=") term)*;
*/
func (p *Parser) comparison() Expression {
	var expr Expression = p.term()

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		var operator Token = p.previous()
		var right Expression = p.term()
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

/*
term -> factor (("+" | "-") factor)*;
*/
func (p *Parser) term() Expression {
	var expr Expression = p.factor()

	for p.match(PLUS, MINUS) {
		var operator Token = p.previous()
		var right Expression = p.factor()
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() Expression {
	var expr Expression = p.unary()

	for p.match(STAR, SLASH, MODULUS) {
		var operator Token = p.previous()
		var right Expression = p.unary()
		expr = &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

/*
unary -> ("!" | "-") unary | primary;
*/
func (p *Parser) unary() Expression {
	for p.match(NOT, MINUS) {
		var operator Token = p.previous()
		var right Expression = p.unary()
		return &Unary{
			operator: operator,
			right:    right,
		}
	}
	return p.primary()
}

/*
primary -> NUMBER | STRING | IDENTIFIER | "true" | "false" | "empty" | "(" expression ")" | incr_decr;
*/
func (p *Parser) primary() Expression {
	if p.match(NUMBER, STRING) {
		return &Literal{
			value: p.previous().literal,
		}
	}

	if p.match(IDENTIFIER) {
		return &Variable{
			name: p.previous(),
		}
	}

	if p.match(TRUE) {
		return &Literal{
			value: true,
		}
	}

	if p.match(FALSE) {
		return &Literal{
			value: false,
		}
	}

	if p.match(EMPTY) {
		return &Literal{
			value: nil,
		}
	}

	// "(" expression ")"
	if p.match(LEFT_PAREN) {
		var expr Expression = p.expression()
		if p.peek().tokenType != RIGHT_PAREN {
			// FIX: throw error here, not return literal
			// ERROR: Expect closing brackets for grouping!
			RuntimeError(p.peek().line, p.peek().lexeme, "expected closing parantheses after statement.")
		} else {
			p.next()
			return &Group{
				expression: expr,
			}
		}
	}

	if p.peek().tokenType == RIGHT_BRACE || p.peek().tokenType == NEWLINE || p.peek().tokenType == EOF {
		return nil
	}
	RuntimeError(p.peek().line, p.peek().lexeme, "unidentified expression.")
	return nil
}

/*
*	Low level functions here
 */
func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if tokenType == p.peek().tokenType {
			p.next()
			return true
		}
	}
	return false
}

func (p *Parser) next() Token {
	var c Token = p.tokens[p.current]
	if !p.end() {
		p.current += 1
	}
	return c
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) peekNext() Token {
	if p.end() {
		return p.tokens[p.current]
	}
	return p.tokens[p.current+1]
}

func (p *Parser) consume(tt TokenType, message string) {
	if p.peek().tokenType == tt {
		p.next()
		return
	}
	panic(message)
}

func (p *Parser) synchronizeToken(endToken TokenType) {
	for !p.end() && p.peek().tokenType != endToken {
		p.next()
	}
	return
}


func (p *Parser) synchronize() {
	// defer func() {
	// 	fmt.Printf("Synchronized at Line %d: %s\n", p.peek().line, p.peek().lexeme)
	// }()
	// consume the token that caused the error
	// p.next()

	for !p.end() {
		if p.match(NEWLINE) {
			/*
				case when synchronizing points to right brace where it checked semicolon exists in previous():
				{
					...
					set a to 1;
				}
			*/
			// if p.peek().tokenType == RIGHT_BRACE {
			// 	p.next()
			// }
			p.next()
			return
		}

		// tokens which usually mark the start of a statement
		// set starting pointer to point to the start of a statement
		tokenType := p.peek().tokenType
		if tokenType == CLASS || tokenType == PROCEDURE || tokenType == SET ||
			tokenType == FOR || tokenType == IF || tokenType == SAY ||
			tokenType == WHILE || tokenType == RETURN || tokenType == EOF || tokenType == LEFT_BRACE {
			fmt.Printf("Synchronized at %s\n", p.peek().lexeme)
			return
		}
		p.next()
	}
}

func (p *Parser) previous() Token {
	if p.current <= 0 {
		return p.tokens[0]
	}
	return p.tokens[p.current-1]
}

func (p *Parser) end() bool {
	return p.tokens[p.current].tokenType == EOF
}

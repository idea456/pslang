package main

/*
Stratified grammar:
expression -> equality;
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

func (p *Parser) Parse() Expression {
	var expr Expression = p.expression()
	return expr
}

func (p *Parser) expression() Expression {
	return p.equality()
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
primary -> NUMBER | STRING | IDENTIFIER | "true" | "false" | "empty" | "(" expression ")";
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
			panic("NOOO there's no closing bracket!")
		} else {
			p.next()
			return &Group{
				expression: expr,
			}
		}
	}
	panic("[Line 1] Unidentified expression!")
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if tokenType == p.tokens[p.current].tokenType {
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

func (p *Parser) previous() Token {
	if p.current <= 0 {
		return p.tokens[0]
	}
	return p.tokens[p.current-1]
}

func (p *Parser) end() bool {
	return p.tokens[p.current].tokenType == EOF
}

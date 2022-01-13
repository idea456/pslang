package main

import "strconv"

type Scanner struct {
	source  string
	tokens  []Token
	current int
	line    int
}

/*
Scanner constructor
*/
func NewScanner(text string) *Scanner {
	scanner := Scanner{}
	scanner.current = 0
	scanner.line = 1
	scanner.source = string(text)
	scanner.tokens = make([]Token, 0)
	return &scanner
}

func (s *Scanner) Scan() []Token {
	for !s.end() {
		s.scanToken()
	}
	// append end
	s.tokens = append(s.tokens, Token{
		tokenType: EOF,
	})
	return s.tokens
}

func (s *Scanner) scanToken() {
	var c string = s.peek()
	if s.match(">", "<", "=", "!") {
		switch c {
		case ">":
			if s.match("=") {
				s.tokens = append(s.tokens, Token{
					tokenType: GREATER_EQUAL,
					literal:   ">=",
					line:      s.line,
				})
			} else {
				s.tokens = append(s.tokens, Token{
					tokenType: GREATER,
					literal:   ">",
					line:      s.line,
				})
			}
		case "<":
			if s.match("=") {
				s.tokens = append(s.tokens, Token{
					tokenType: LESS_EQUAL,
					literal:   "<=",
					line:      s.line,
				})
			} else {
				s.tokens = append(s.tokens, Token{
					tokenType: LESS,
					literal:   "<",
					line:      s.line,
				})
			}
		case "!":
			if s.match("=") {
				s.tokens = append(s.tokens, Token{
					tokenType: NOT_EQUAL,
					literal:   "!=",
					line:      s.line,
				})
			} else {
				s.tokens = append(s.tokens, Token{
					tokenType: NOT,
					literal:   "!",
					line:      s.line,
				})
			}
		case "=":
			if s.match("=") {
				s.tokens = append(s.tokens, Token{
					tokenType: EQUAL_EQUAL,
					literal:   "==",
					line:      s.line,
				})
			}
		}

	}
	// scan for operators, brackets and semicolon
	if s.match("+", "-", "*", "/", "%", "(", ")", "{", "}", ";", ",", ".") {
		s.tokens = append(s.tokens, Token{
			tokenType: keywords[s.previous()],
			lexeme:    c,
			line:      s.line,
		})
	}

	// ignore whitespaces, tabs and newlines
	if s.match("\n", "\r", " ", "\t") {
		if s.match("\n") {
			s.line += 1
		}
		return
	}

	// scan for numbers
	if s.isNumber(c) {
		s.scanNumber()
	}

	if s.match("\"") {
		s.scanString()
	}

	// scan for keywords or variable names
	if s.isAlpha(c) {
		s.scanIdentifier()
	}

	// error: unexpected character
	// raise Error
}

func (s *Scanner) scanNumber() {
	var numStr string = ""
	for s.peek() >= "0" && s.peek() <= "9" && !s.end() {
		numStr += s.next()
		// consume decimal point also
		if s.peek() == "." {
			numStr += s.next()
		}
	}

	num, _ := strconv.Atoi(numStr)

	s.tokens = append(s.tokens, Token{
		tokenType: NUMBER,
		lexeme:    numStr,
		literal:   num,
		line:      s.line,
	})
}

func (s *Scanner) scanString() {
	var value string = ""
	for s.peek() != "\"" && !s.end() {
		value += s.next()
	}

	if s.end() && s.peek() != "'" {
		// error: unterminated string
		return
	}

	// consume closing "
	s.next()

	s.tokens = append(s.tokens, Token{
		tokenType: STRING,
		lexeme:    value,
		literal:   value,
		line:      s.line,
	})
}

func (s *Scanner) scanIdentifier() {
	var name string = ""
	// case when e.g. or vs order, perform maximal munching
	for (s.isAlpha(s.peek()) || s.isNumber((s.peek()))) && !s.end() {
		name += s.next()
	}

	keyword, exist := keywords[name]
	if !exist {
		s.tokens = append(s.tokens, Token{
			tokenType: IDENTIFIER,
			lexeme:    name,
			literal:   name,
			line:      s.line,
		})
	} else {
		s.tokens = append(s.tokens, Token{
			tokenType: keyword,
			lexeme:    name,
			literal:   name,
			line:      s.line,
		})
	}

}

/*
@return
next() only mutates the current pointer and returns the current character being consumed
*/
func (s *Scanner) next() string {
	var c byte = s.source[s.current]
	s.current += 1
	return string(c)
}

func (s *Scanner) peek() string {
	if s.end() {
		return string(s.source[len(s.source)-1])
	}
	return string(s.source[s.current])
}

func (s *Scanner) previous() string {
	if s.current <= 0 {
		return string(s.source[0])
	}
	return string(s.source[s.current-1])
}

func (s *Scanner) isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z")
}

func (s *Scanner) isNumber(c string) bool {
	return c >= "0" && c <= "9"
}

func (s *Scanner) end() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) match(keywords ...string) bool {
	for _, kw := range keywords {
		if s.peek() == kw {
			s.next()
			return true
		}
	}
	return false
}

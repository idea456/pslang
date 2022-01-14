package main

import "fmt"

type RuntimeError struct {
	line    int
	token   Token
	message string
}

func (e *RuntimeError) Error() {
	fmt.Sprintf("[Line %d] Runtime error at '%s' : %s", e.line, e.token.lexeme, e.message)
}

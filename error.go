package main

import "fmt"

func RuntimeError(line int, lexeme interface{}, message string) {
	lexemeStr := fmt.Sprint(lexeme)
	var msg string = fmt.Sprintf("[Line %d] Runtime Error at '%s': %s", line, lexemeStr, message)
	panic(msg)
}

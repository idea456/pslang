package main

import "fmt"

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	var itpr Interpreter = Interpreter{}
	itpr.environment = NewEnvironment()
	return &itpr
}

func (itpr *Interpreter) Interpret(stmts []Statement) {
	for _, stmt := range stmts {
		itpr.execute(stmt)
	}
}

// func (itpr *Interpreter) accept(visitor VisitorStmt) {
// 	return
// }

// func (itpr *Interpreter) accept(visitor VisitorExpr) interface{} {
// 	return 0
// }

func (itpr *Interpreter) evaluate(expr Expression) interface{} {
	return expr.accept(itpr)
}

// execution for statements
func (itpr *Interpreter) execute(stmt Statement) {
	stmt.accept(itpr)
}

func (itpr *Interpreter) visitBinaryExpr(expr *Binary) interface{} {
	var left interface{} = itpr.evaluate(expr.left)
	var right interface{} = itpr.evaluate(expr.right)

	switch expr.operator.tokenType {
	case PLUS:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) + itpr.toString(right)
		}
		return itpr.toNum(left) + itpr.toNum(right)
	case MINUS:
		return itpr.toNum(left) - itpr.toNum(right)
	case STAR:
		return itpr.toNum(left) * itpr.toNum(right)
	case SLASH:
		if itpr.toNum(right) == 0 {
			panic("NOOO cannot divide by 0!")
		}
		return itpr.toNum(left) / itpr.toNum(right)
	case EQUAL_EQUAL:
		if left == nil || right == nil {
			return false
		}
		return left == right
	case GREATER:
		// comparisons are only supported between strings and integers
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) > itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) > itpr.toNum(right)
		}
		panic("Error, expected string or integer for comparisons!")
	case GREATER_EQUAL:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) >= itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) >= itpr.toNum(right)
		}
		panic("Error, expected string or integer for comparisons!")
	case LESS:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) < itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) < itpr.toNum(right)
		}
		panic("Error, expected string or integer for comparisons!")
	case LESS_EQUAL:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) <= itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) <= itpr.toNum(right)
		}
		panic("Error, expected string or integer for comparisons!")
	}
	return nil
}

func (itpr *Interpreter) visitLiteralExpr(expr *Literal) interface{} {
	return expr.value
}

func (itpr *Interpreter) visitUnaryExpr(expr *Unary) interface{} {
	var right interface{} = itpr.evaluate(expr.right)
	switch expr.operator.tokenType {
	case MINUS:
		return -itpr.toNum(right)
	case NOT:
		return !itpr.evaluateBool(right)
	}
	return nil
}

func (itpr *Interpreter) visitGroupExpr(expr *Group) interface{} {
	return itpr.evaluate(expr.expression)
}

func (itpr *Interpreter) visitVariableExpr(expr *Variable) interface{} {
	return (*itpr.environment).Get(expr.name)
}

func (itpr *Interpreter) visitVariableStmt(stmt *VariableStmt) {
	var value interface{} = itpr.evaluate(stmt.initializer)
	(*itpr.environment).Set(stmt.name, value)
}

func (itpr *Interpreter) visitSayStmt(stmt *SayStmt) {
	var value interface{} = itpr.evaluate(stmt.expression)
	fmt.Println(value)
}

func (itpr *Interpreter) visitBlockStmt(stmt *BlockStmt) {
	fmt.Println("block")
}

func (itpr *Interpreter) visitExprStmt(stmt *ExprStmt) {
	itpr.evaluate(stmt.expression)
}

func (itpr *Interpreter) evaluateBool(expr interface{}) bool {
	if expr == nil {
		return false
	}

	truth, ok := expr.(bool)
	if ok {
		return truth
	}
	// assume all other values are true
	return true
}

// type SignedNum interface {
// 	~int | ~int8 | ~int16 | ~int32 | ~int64
// }

func (itpr *Interpreter) toString(expr interface{}) string {
	text, ok := expr.(string)
	if !ok {
		panic("Error, string expected!")
	}
	return text
}

func (itpr *Interpreter) isString(expr interface{}) bool {
	_, ok := expr.(string)
	return ok
}

func (itpr *Interpreter) isNum(expr interface{}) bool {
	switch expr.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

func (itpr *Interpreter) toNum(expr interface{}) float64 {
	switch t := expr.(type) {
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return float64(t)
	default:
		panic("Error, integer expected!")
	}
}

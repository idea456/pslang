package main

type Interpreter struct {
}

func (itpr *Interpreter) evaluate(expr Expression) interface{} {
	return expr.accept(itpr)
}

func (itpr *Interpreter) visitBinaryExpr(expr *Binary) interface{} {
	var left interface{} = itpr.evaluate(expr.left)
	var right interface{} = itpr.evaluate(expr.right)

	switch expr.operator.tokenType {
	case PLUS:
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
	}
	return left
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
	return 0
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

func (itpr *Interpreter) toNum(expr interface{}) int {
	// use type assertion
	num, ok := expr.(int)
	if !ok {
		// ERROR: should be num
		panic("Error, integer expected!")
	}
	return num
}

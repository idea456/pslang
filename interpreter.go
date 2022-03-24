package main

import "fmt"

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	var itpr Interpreter = Interpreter{}
	itpr.environment = NewEnv()
	return &itpr
}

func (itpr *Interpreter) Interpret(stmts []Statement) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%+v\n", r)
		}
	}()

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

func (itpr *Interpreter) visitLogicalExpr(expr *Logical) interface{} {
	var left interface{} = itpr.evaluate(expr.left)
	if expr.operator.tokenType == OR {
		if itpr.evaluateBool(left) {
			return left
		}
	} else {
		if !itpr.evaluateBool(left) {
			return left
		}
	}
	return itpr.evaluate(expr.right)
}

func (itpr *Interpreter) visitBinaryExpr(expr *Binary) interface{} {
	var left interface{} = itpr.evaluate(expr.left)
	var right interface{} = itpr.evaluate(expr.right)

	checkedComparison := false
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
			RuntimeError(expr.operator.line, right, "cannot divide numbers by 0.")
		}
		return itpr.toNum(left) / itpr.toNum(right)
	case MODULUS:
		leftNum, okLeft := left.(int)
		rightNum, okRight := right.(int)
		if !okLeft || !okRight {
			RuntimeError(expr.operator.line, right, "cannot modulus non-integers!")
		}
		return leftNum % rightNum
	case EQUAL_EQUAL:
		if left == nil || right == nil {
			return false
		}
		return left == right
	case NOT_EQUAL:
		if left == nil || right == nil {
			return false
		}
		return left != right
	case GREATER:
		// comparisons are only supported between strings and integers
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) > itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) > itpr.toNum(right)
		}
		checkedComparison = true
	case GREATER_EQUAL:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) >= itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) >= itpr.toNum(right)
		}
		checkedComparison = true
	case LESS:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) < itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) < itpr.toNum(right)
		}
		checkedComparison = true
	case LESS_EQUAL:
		if itpr.isString(left) && itpr.isString(right) {
			return itpr.toString(left) <= itpr.toString(right)
		}
		if itpr.isNum(left) && itpr.isNum(right) {
			return itpr.toNum(left) <= itpr.toNum(right)
		}
		checkedComparison = true
	}
	if checkedComparison {
		RuntimeError(expr.operator.line, expr.operator.lexeme, "Error, expected string or integer for comparisons!")
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

func (itpr *Interpreter) visitBlockStmt(blockStmt *BlockStmt) {
	var enclosing *Environment = itpr.environment
	defer func() {
		itpr.environment = enclosing
	}()

	itpr.environment = NewEnv()
	itpr.environment.enclosing = enclosing
	for _, stmt := range blockStmt.statements {
		itpr.execute(stmt)
	}
}

func (itpr *Interpreter) visitIfStmt(stmt *IfStmt) {
	if itpr.evaluateBool(itpr.evaluate(stmt.expression)) {
		itpr.execute(stmt.thenBranch)
	} else {
		itpr.execute(stmt.elseBranch)
	}
}

func (itpr *Interpreter) visitWhileStmt(stmt *WhileStmt) {
	var condition interface{} = itpr.evaluate(stmt.condition)

	for itpr.evaluateBool(condition) {
		itpr.execute(stmt.body)
		// re-evaluate the condition again after executing a statement in the body
		condition = itpr.evaluate(stmt.condition)
	}
}

func (itpr *Interpreter) visitExprStmt(stmt *ExprStmt) {
	itpr.evaluate(stmt.expression)
}

func (itpr *Interpreter) visitIncrDecrStmt(stmt *IncrDecrStmt) {
	var left interface{} = (*itpr.environment).Get(stmt.identifier)
	var right interface{} = itpr.evaluate(stmt.right)

	if !(itpr.isNum(left) && itpr.isNum(right)) {
		RuntimeError(stmt.identifier.line, stmt.identifier.lexeme, "only numbers allowed for increments/decrements.")
	}

	if stmt.operator.tokenType == INCREMENT {
		(*itpr.environment).Set(stmt.identifier, itpr.toNum(left)+itpr.toNum(right))
	} else if stmt.operator.tokenType == DECREMENT {
		(*itpr.environment).Set(stmt.identifier, itpr.toNum(left)-itpr.toNum(right))
	}
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

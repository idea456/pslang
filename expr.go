package main

type VisitorExpr interface {
	visitLiteralExpr(*Literal) interface{}
	visitUnaryExpr(*Unary) interface{}
	visitBinaryExpr(*Binary) interface{}
	visitVariableExpr(*Variable) interface{}
	visitGroupExpr(*Group) interface{}
}

type Expression interface {
	accept(VisitorExpr) interface{}
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor VisitorExpr) interface{} {
	return visitor.visitLiteralExpr(expr)
}

type Unary struct {
	operator Token
	right    Expression
}

func (expr *Unary) accept(visitor VisitorExpr) interface{} {
	return visitor.visitUnaryExpr(expr)
}

type Binary struct {
	left     Expression
	operator Token
	right    Expression
}

func (expr *Binary) accept(visitor VisitorExpr) interface{} {
	return visitor.visitBinaryExpr(expr)
}

type Variable struct {
	name Token
}

func (expr *Variable) accept(visitor VisitorExpr) interface{} {
	return visitor.visitVariableExpr(expr)
}

type Group struct {
	expression Expression
}

func (expr *Group) accept(visitor VisitorExpr) interface{} {
	return visitor.visitGroupExpr(expr)
}

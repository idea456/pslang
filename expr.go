package main

type Visitor interface {
	visitLiteralExpr(*Literal)
	visitUnaryExpr(*Unary)
	visitBinaryExpr(*Binary)
	visitVariableExpr(*Variable)
	visitGroupExpr(*Group)
}

type Expression interface {
	accept(Visitor)
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(visitor Visitor) {
	visitor.visitLiteralExpr(expr)
}

type Unary struct {
	operator Token
	right    Expression
}

func (expr *Unary) accept(visitor Visitor) {
	visitor.visitUnaryExpr(expr)
}

type Binary struct {
	left     Expression
	operator Token
	right    Expression
}

func (expr *Binary) accept(visitor Visitor) {
	visitor.visitBinaryExpr(expr)
}

type Variable struct {
	name Token
}

func (expr *Variable) accept(visitor Visitor) {
	visitor.visitVariableExpr(expr)
}

type Group struct {
	expression Expression
}

func (expr *Group) accept(visitor Visitor) {
	visitor.visitGroupExpr(expr)
}

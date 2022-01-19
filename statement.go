package main

type VisitorStmt interface {
	visitVariableStmt(stmt *VariableStmt)
	visitSayStmt(stmt *SayStmt)
	visitBlockStmt(stmt *BlockStmt)
	visitExprStmt(stmt *ExprStmt)
	visitIncrDecrStmt(stmt *IncrDecrStmt)
	visitIfStmt(stmt *IfStmt)
}

type Statement interface {
	accept(VisitorStmt)
}

type VariableStmt struct {
	name        Token
	initializer Expression
}

func (stmt *VariableStmt) accept(visitor VisitorStmt) {
	visitor.visitVariableStmt(stmt)
}

type SayStmt struct {
	expression Expression
}

func (stmt *SayStmt) accept(visitor VisitorStmt) {
	visitor.visitSayStmt(stmt)
}

type BlockStmt struct {
	statements []Statement
}

func (stmt *BlockStmt) accept(visitor VisitorStmt) {
	visitor.visitBlockStmt(stmt)
}

type ExprStmt struct {
	expression Expression
}

func (stmt *ExprStmt) accept(visitor VisitorStmt) {
	visitor.visitExprStmt(stmt)
}

type IncrDecrStmt struct {
	identifier Token
	operator   Token
	right      Expression
}

func (stmt *IncrDecrStmt) accept(visitor VisitorStmt) {
	visitor.visitIncrDecrStmt(stmt)
}

type IfStmt struct {
	expression Expression
	thenBranch Statement
	elseBranch Statement
}

func (stmt *IfStmt) accept(visitor VisitorStmt) {
	visitor.visitIfStmt(stmt)
}

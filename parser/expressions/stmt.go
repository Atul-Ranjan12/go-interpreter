package expressions

import "github.com/Atul-Ranjan12/token"

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) (interface{}, error)
	VisitExprStatementStmt(stmt *ExprStatement) (interface{}, error)
	VisitPrintStatementStmt(stmt *PrintStatement) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
}

// These are functions for Block 
type Block struct {
	Statements []Stmt
}

var _ Stmt = (*Block)(nil)

func (e *Block) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitBlockStmt(e)
}

// These are functions for ExprStatement 
type ExprStatement struct {
	Expression Expr
}

var _ Stmt = (*ExprStatement)(nil)

func (e *ExprStatement) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExprStatementStmt(e)
}

// These are functions for PrintStatement 
type PrintStatement struct {
	Expression Expr
}

var _ Stmt = (*PrintStatement)(nil)

func (e *PrintStatement) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStatementStmt(e)
}

// These are functions for Var 
type Var struct {
	Name token.Token
	Initializer Expr
}

var _ Stmt = (*Var)(nil)

func (e *Var) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitVarStmt(e)
}


package expressions

import "github.com/Atul-Ranjan12/token"

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) (interface{}, error)
	VisitExprStatementStmt(stmt *ExprStatement) (interface{}, error)
	VisitPrintStatementStmt(stmt *PrintStatement) (interface{}, error)
	VisitReturnStmt(stmt *Return) (interface{}, error)
	VisitWhileStatementStmt(stmt *WhileStatement) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
	VisitIfStmt(stmt *If) (interface{}, error)
	VisitFunctionStmt(stmt *Function) (interface{}, error)
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

// These are functions for Return 
type Return struct {
	Keyword token.Token
	Value Expr
}

var _ Stmt = (*Return)(nil)

func (e *Return) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitReturnStmt(e)
}

// These are functions for WhileStatement 
type WhileStatement struct {
	Condition Expr
	Body Stmt
}

var _ Stmt = (*WhileStatement)(nil)

func (e *WhileStatement) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitWhileStatementStmt(e)
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

// These are functions for If 
type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

var _ Stmt = (*If)(nil)

func (e *If) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitIfStmt(e)
}

// These are functions for Function 
type Function struct {
	Name token.Token
	Params []token.Token
	Body []Stmt
}

var _ Stmt = (*Function)(nil)

func (e *Function) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitFunctionStmt(e)
}


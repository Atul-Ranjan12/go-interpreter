package expressions

type Stmt interface {
	Accept(visitor StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) (interface{}, error)
	VisitPrintStmt(stmt *Print) (interface{}, error)
}

// These are functions for Expression
type Expression struct {
	Expression Expr
}

var _ Stmt = (*Expression)(nil)

func (e *Expression) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitExpressionStmt(e)
}

// These are functions for Print
type Print struct {
	Expression Expr
}

var _ Stmt = (*Print)(nil)

func (e *Print) Accept(visitor StmtVisitor) (interface{}, error) {
	return visitor.VisitPrintStmt(e)
}

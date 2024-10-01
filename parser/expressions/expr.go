package expressions

import "github.com/Atul-Ranjan12/token"

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
}

// These are functions for Assign 
type Assign struct {
	Name token.Token
	Value Expr
}

var _ Expr = (*Assign)(nil)

func (e *Assign) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitAssignExpr(e)
}

// These are functions for Binary 
type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

var _ Expr = (*Binary)(nil)

func (e *Binary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(e)
}

// These are functions for Grouping 
type Grouping struct {
	Expression Expr
}

var _ Expr = (*Grouping)(nil)

func (e *Grouping) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(e)
}

// These are functions for Literal 
type Literal struct {
	Value interface{}
}

var _ Expr = (*Literal)(nil)

func (e *Literal) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(e)
}

// These are functions for Unary 
type Unary struct {
	Operator token.Token
	Right Expr
}

var _ Expr = (*Unary)(nil)

func (e *Unary) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(e)
}

// These are functions for Variable 
type Variable struct {
	Name token.Token
}

var _ Expr = (*Variable)(nil)

func (e *Variable) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitVariableExpr(e)
}


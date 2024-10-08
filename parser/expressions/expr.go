package expressions

import "github.com/Atul-Ranjan12/token"

type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGetExpr(expr *Get) (interface{}, error)
	VisitSetExpr(expr *Set) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitBreakExprExpr(expr *BreakExpr) (interface{}, error)
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

// These are functions for Logical 
type Logical struct {
	Left Expr
	Right Expr
	Operator token.Token
}

var _ Expr = (*Logical)(nil)

func (e *Logical) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLogicalExpr(e)
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

// These are functions for Call 
type Call struct {
	Callee Expr
	Paren token.Token
	Arguments []Expr
}

var _ Expr = (*Call)(nil)

func (e *Call) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitCallExpr(e)
}

// These are functions for Get 
type Get struct {
	Object Expr
	Name token.Token
}

var _ Expr = (*Get)(nil)

func (e *Get) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGetExpr(e)
}

// These are functions for Set 
type Set struct {
	Object Expr
	Name token.Token
	Value Expr
}

var _ Expr = (*Set)(nil)

func (e *Set) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitSetExpr(e)
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

// These are functions for BreakExpr 
type BreakExpr struct {
}

var _ Expr = (*BreakExpr)(nil)

func (e *BreakExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBreakExprExpr(e)
}


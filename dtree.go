package dtree

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/Knetic/govaluate"
)

var (
	ErrUndecidable      = errors.New("undecidable")
	ErrInvalidOutcome   = errors.New("invalid outcome")
	ErrInvalidCondition = errors.New("invalid condition")
)

type Tree struct {
	Root Node
}

type Node interface {
	Next(params map[string]interface{}) (Node, error)
}

func (t *Tree) Decide(params map[string]interface{}) (interface{}, error) {
	var err error
	node := t.Root
	for _, ok := node.(*Condition); ok; _, ok = node.(*Condition) {
		node, err = node.Next(params)
		if err != nil {
			return nil, err
		}
	}
	outcome, ok := node.(*Outcome)
	if !ok {
		return nil, ErrUndecidable
	}
	return outcome.Value, nil
}

type Outcome struct {
	Value interface{}
}

func (o *Outcome) Next(_ map[string]interface{}) (Node, error) {
	return nil, nil
}

// NewOutcome returns a new Outcome node.
// The value must be integer literal, float literal, string literal (with or without "").
func NewOutcome(value interface{}) (*Outcome, error) {
	expr, err := parser.ParseExpr(fmt.Sprint(value))
	if err != nil {
		return nil, ErrInvalidOutcome
	}
	switch expr.(type) {
	case *ast.Ident, *ast.BasicLit:
		return &Outcome{Value: value}, nil
	default:
		return nil, ErrInvalidOutcome
	}
}

type Condition struct {
	EvaluablePredicate *govaluate.EvaluableExpression
	Branches           map[interface{}]Node
}

// NewCondition returns a new Condition node. The value must be binary expression.
// If value accepts unary expression, it'd be ambiguous for a value=`X`. Because the value
// can be a boolean expression `X` (equivalent to `X == true`) or a string outcome `X`.
func NewCondition(predicate string) (*Condition, error) {
	expr, err := parser.ParseExpr(predicate)
	if err != nil {
		return nil, ErrInvalidCondition
	}
	if _, ok := expr.(*ast.BinaryExpr); !ok {
		// Condition should not be unary expression to distinguish condition and outcome.
		return nil, ErrInvalidCondition
	}
	evaluablePredicate, err := govaluate.NewEvaluableExpression(predicate)
	if err != nil {
		return nil, ErrInvalidCondition
	}
	return &Condition{
		EvaluablePredicate: evaluablePredicate,
		Branches:           map[interface{}]Node{},
	}, nil
}

func (c *Condition) Next(params map[string]interface{}) (Node, error) {
	value, err := c.EvaluablePredicate.Evaluate(params)
	if err != nil {
		return nil, err
	}
	node, ok := c.Branches[value]
	if !ok {
		return nil, ErrUndecidable
	}
	return node, nil
}

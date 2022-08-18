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
	Root *Node
}

func (t *Tree) Decide(params map[string]interface{}) (interface{}, error) {
	var err error
	node := t.Root
	for node != nil && node.IsCondition() {
		node, err = node.Condition.Decide(params)
		if err != nil {
			return nil, err
		}
	}
	if node == nil {
		return ErrUndecidable, nil
	}
	return node.Outcome.Value, nil
}

type Node struct {
	Outcome   *Outcome
	Condition *Condition
}

func NewConditionNode(condition *Condition) *Node {
	return &Node{Condition: condition}
}

func NewOutcomeNode(outcome *Outcome) *Node {
	return &Node{Outcome: outcome}
}

func (n *Node) IsCondition() bool {
	return n.Condition != nil
}

type Outcome struct {
	Value interface{}
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
	Predicate *govaluate.EvaluableExpression
	Branches  map[interface{}]*Node
}

// NewCondition returns a new Condition node. The value must be binary expression.
// If value accepts unary expression, it'd be ambiguous for a value=`X`.
// It can be a boolean expression `X` (equivalent to `X == true`) or a string outcome `X`.
func NewCondition(value string) (*Condition, error) {
	expr, err := parser.ParseExpr(value)
	if err != nil {
		return nil, ErrInvalidCondition
	}
	if _, ok := expr.(*ast.BinaryExpr); !ok {
		// Condition should not be unary expression to distinguish condition and outcome.
		return nil, ErrInvalidCondition
	}
	predicate, err := govaluate.NewEvaluableExpression(value)
	if err != nil {
		return nil, ErrInvalidCondition
	}
	return &Condition{Predicate: predicate, Branches: map[interface{}]*Node{}}, nil
}

func (c *Condition) Decide(params map[string]interface{}) (*Node, error) {
	value, err := c.Predicate.Evaluate(params)
	if err != nil {
		return nil, err
	}
	node, ok := c.Branches[value]
	if !ok {
		return nil, ErrUndecidable
	}
	return node, nil
}

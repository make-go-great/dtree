package dtree

import (
	"encoding/json"
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
	Root *Node `json:"root"`
}

func NewTreeFromJson(bytes []byte) (*Tree, error) {
	tree := Tree{}
	if err := json.Unmarshal(bytes, &tree); err != nil {
		return nil, err
	}
	if err := tree.Initialize(); err != nil {
		return nil, err
	}
	return &tree, nil
}

func (t *Tree) Initialize() error {
	if t == nil || t.Root == nil || t.Root.Condition == nil {
		return nil
	}

	queue := []*Condition{t.Root.Condition}
	for len(queue) > 0 {
		head := queue[0]
		queue = queue[1:]
		if err := head.Initialize(); err != nil {
			return err
		}

		for _, branch := range head.Branches {
			if branch.Next != nil && branch.Next.Condition != nil {
				queue = append(queue, branch.Next.Condition)
			}
		}
	}

	return nil
}

type Branch struct {
	Value interface{} `json:"value"`
	Next  *Node       `json:"next"`
}

type Node struct {
	*Outcome
	*Condition
}

func (t *Tree) Decide(params map[string]interface{}) (interface{}, error) {
	var err error
	node := t.Root
	for node != nil && node.Condition != nil {
		node, err = node.Condition.Next(params)
		if err != nil {
			return nil, err
		}
	}
	if node == nil || node.Outcome == nil {
		return nil, ErrUndecidable
	}
	return node.Outcome.Value, nil
}

type Outcome struct {
	Value interface{} `json:"outcome"`
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
	valueToNextNode    map[interface{}]*Node
	evaluablePredicate *govaluate.EvaluableExpression

	Predicate string    `json:"predicate"`
	Branches  []*Branch `json:"branches"`
}

func (c *Condition) Initialize() error {
	evaluablePredicate, err := govaluate.NewEvaluableExpression(c.Predicate)
	if err != nil {
		return err
	}
	c.evaluablePredicate = evaluablePredicate

	c.valueToNextNode = make(map[interface{}]*Node, len(c.Branches))
	for _, branch := range c.Branches {
		c.valueToNextNode[branch.Value] = branch.Next
	}

	return nil
}

func (c *Condition) AddBranch(value interface{}, nextNode *Node) {
	if c.valueToNextNode == nil {
		c.valueToNextNode = make(map[interface{}]*Node)
	}
	c.valueToNextNode[value] = nextNode

	branch := &Branch{
		Value: value,
		Next:  nextNode,
	}
	c.Branches = append(c.Branches, branch)
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
		Predicate:          predicate,
		evaluablePredicate: evaluablePredicate,
		valueToNextNode:    map[interface{}]*Node{},
	}, nil
}

func (c *Condition) Next(params map[string]interface{}) (*Node, error) {
	value, err := c.evaluablePredicate.Evaluate(params)
	if err != nil {
		return nil, err
	}
	node, ok := c.valueToNextNode[value]
	if !ok {
		return nil, ErrUndecidable
	}
	return node, nil
}

package solver

import (
	"fmt"

	"github.com/Paul-weqe/orac-verify/internal/encoder"
	z3 "github.com/Z3Prover/z3/src/api/go"
)

// Result is the outcome of a verification check.
type Result struct {
	Satisfiable bool              // true if the negated rule has a model
	Model       map[string]string // variable name to value, if satisfiable
}

// Check translates the assertion tree into Z3, solves it, and returns the
// result. A satisfiable result means the negated rule holds, so the rule
// was violated and Model names the offending values. An unsatisfiable
// result means the rule holds for the input data.
func Check(a encoder.Assertion) (Result, error) {
	ctx := z3.NewContext()
	s := ctx.NewSolver()

	expr, err := translate(ctx, a)
	if err != nil {
		return Result{}, err
	}
	s.Assert(expr)

	switch s.Check() {
	case z3.Unsatisfiable:
		return Result{Satisfiable: false}, nil
	case z3.Satisfiable:
		model := s.Model()
		values := extractModel(model)
		return Result{Satisfiable: true, Model: values}, nil
	default:
		return Result{}, fmt.Errorf("solver returned unknown")
	}
}

// translate walks the assertion tree and produces a Z3 boolean expression.
func translate(ctx *z3.Context, a encoder.Assertion) (*z3.Expr, error) {
	switch n := a.(type) {
	case encoder.BoolLit:
		if n.Value {
			return ctx.MkTrue(), nil
		}
		return ctx.MkFalse(), nil

	case encoder.IntCompare:
		left, err := translateExpr(ctx, n.Left)
		if err != nil {
			return nil, err
		}
		right, err := translateExpr(ctx, n.Right)
		if err != nil {
			return nil, err
		}
		switch n.Op {
		case encoder.OpEq:
			return ctx.MkEq(left, right), nil
		case encoder.OpNe:
			return ctx.MkNot(ctx.MkEq(left, right)), nil
		case encoder.OpLt:
			return ctx.MkLt(left, right), nil
		case encoder.OpLe:
			return ctx.MkLe(left, right), nil
		case encoder.OpGt:
			return ctx.MkGt(left, right), nil
		case encoder.OpGe:
			return ctx.MkGe(left, right), nil
		default:
			return nil, fmt.Errorf("unknown compare op: %d", n.Op)
		}

	case encoder.And:
		if len(n.Children) == 0 {
			return ctx.MkTrue(), nil
		}
		children := make([]*z3.Expr, 0, len(n.Children))
		for _, child := range n.Children {
			c, err := translate(ctx, child)
			if err != nil {
				return nil, err
			}
			children = append(children, c)
		}
		if len(children) == 1 {
			return children[0], nil
		}
		return ctx.MkAnd(children...), nil

	case encoder.Or:
		if len(n.Children) == 0 {
			return ctx.MkFalse(), nil
		}
		children := make([]*z3.Expr, 0, len(n.Children))
		for _, child := range n.Children {
			c, err := translate(ctx, child)
			if err != nil {
				return nil, err
			}
			children = append(children, c)
		}
		if len(children) == 1 {
			return children[0], nil
		}
		return ctx.MkOr(children...), nil

	case encoder.Not:
		child, err := translate(ctx, n.Children)
		if err != nil {
			return nil, err
		}
		return ctx.MkNot(child), nil

	default:
		return nil, fmt.Errorf("unknown assertion type: %T", a)
	}
}

// translateExpr walks an integer expression and produces a Z3 integer
// expression.
func translateExpr(ctx *z3.Context, e encoder.Expr) (*z3.Expr, error) {
	switch n := e.(type) {
	case encoder.IntVar:
		return ctx.MkIntConst(n.Name), nil
	case encoder.IntLit:
		return ctx.MkInt(int(n.Value), ctx.MkIntSort()), nil
	default:
		return nil, fmt.Errorf("unknown expr type: %T", e)
	}
}

// extractModel reads every declared constant from the model and returns
// its value as a string.
func extractModel(model *z3.Model) map[string]string {
	values := make(map[string]string)
	for i := uint(0); i < model.NumConsts(); i++ {
		decl := model.GetConstDecl(i)
		name := decl.GetName().String()
		val := model.GetConstInterp(decl)
		if val != nil {
			values[name] = val.String()
		}
	}
	return values
}

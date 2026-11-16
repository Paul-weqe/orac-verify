package encoder

type CompareOp int

const (
	OpEq CompareOp = iota
	OpNe
	OpLt
	OpLe
	OpGt
	OpGe
)

type Assertion interface {
	isAssertion()
}

type IntCompare struct {
	Left  Expr
	Op    CompareOp
	Right Expr
}

func (IntCompare) isAssertion() {}

type And struct {
	Children []Assertion
}

func (And) isAssertion() {}

type Or struct {
	Children []Assertion
}

func (Or) isAssertion() {}

type Not struct {
	Children Assertion
}

func (Not) isAssertion() {}

type BoolLit struct {
	Value bool
}

func (BoolLit) isAssertion() {}

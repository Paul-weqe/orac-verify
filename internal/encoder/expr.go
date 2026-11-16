package encoder

type Expr interface {
	isExpr()
}

type IntVar struct {
	Name string
}

func (IntVar) isExpr() {}

type IntLit struct {
	Value int64
}

func (IntLit) isExpr() {}

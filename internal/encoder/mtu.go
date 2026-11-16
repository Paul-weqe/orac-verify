package encoder

import (
	"fmt"

	"github.com/Paul-weqe/orac-data/models"
)

type MTUConstraint struct {
	Min int
}

func EncodeMTU(data *models.InterfaceData, c MTUConstraint) (Assertion, error) {
	if data == nil {
		return nil, fmt.Errorf("nil interface data")
	}

	if c.Min < 0 {
		return nil, fmt.Errorf("MTU constraint: Min must be non-negative, got %d", c.Min)
	}

	if len(data.Interfaces) == 0 {
		return BoolLit{Value: false}, nil
	}

	facts := make([]Assertion, 0, len(data.Interfaces))
	violations := make([]Assertion, 0, len(data.Interfaces))

	for _, iface := range data.Interfaces {
		name := mtuVarName(iface.Name)
		facts = append(facts, IntCompare{
			Left:  IntVar{Name: name},
			Op:    OpEq,
			Right: IntLit{Value: int64(iface.MTU)},
		})

		violations = append(violations, IntCompare{
			Left:  IntVar{Name: name},
			Op:    OpLt,
			Right: IntLit{Value: int64(c.Min)},
		})
	}

	var violation Assertion
	if len(violations) == 1 {
		violation = violations[0]
	} else {
		violation = Or{Children: violations}
	}

	var factTree Assertion
	if len(facts) == 1 {
		factTree = violations[0]
	} else {
		factTree = And{Children: facts}
	}

	return And{Children: []Assertion{factTree, violation}}, nil
}

func mtuVarName(ifaceName string) string {
	return "mtu_" + ifaceName
}

package expression

import (
	"fmt"
)

type Expression interface {
	Expression()
	Equals(other Expression) bool
}

type Not struct {
	Exp Expression
}

func (f *Not) Expression() {}

func (f *Not) Equals(other Expression) bool {
	if casted, ok := other.(*Not); ok {
		return f.Exp.Equals(casted.Exp)
	} else {
		return false
	}
}

func (f *Not) String() string {
	return fmt.Sprintf("NOT{%v}", f.Exp)
}

type Ors struct {
	Exps []Expression
}

func (f *Ors) Expression() {}

func (f *Ors) Equals(other Expression) bool {
	if casted, ok := other.(*Ors); ok {
		if len(f.Exps) != len(casted.Exps) {
			return false
		}
		for i := range f.Exps {
			if !f.Exps[i].Equals(casted.Exps[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (f *Ors) String() string {
	return fmt.Sprintf("OR{%v}", f.Exps)
}

type Ands struct {
	Exps []Expression
}

func (f *Ands) Expression() {}

func (f *Ands) Equals(other Expression) bool {
	if casted, ok := other.(*Ands); ok {
		if len(f.Exps) != len(casted.Exps) {
			return false
		}
		for i := range f.Exps {
			if !f.Exps[i].Equals(casted.Exps[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (f *Ands) String() string {
	return fmt.Sprintf("AND{%v}", f.Exps)
}

type UnaryExpression struct {
	Field string
	Value string
}

func (f *UnaryExpression) Expression() {}

func (f *UnaryExpression) Equals(other Expression) bool {
	if casted, ok := other.(*UnaryExpression); ok {
		return f.Field == casted.Field && f.Value == casted.Value
	} else {
		return false
	}
}

func (f *UnaryExpression) String() string {
	return fmt.Sprintf("%v IS %v", f.Field, f.Value)
}

type BinaryExpression struct {
	Field string
	Op    string
	Value string
}

func (f *BinaryExpression) Expression() {}

func (f *BinaryExpression) Equals(other Expression) bool {
	if casted, ok := other.(*BinaryExpression); ok {
		return f.Field == casted.Field && f.Op == casted.Op && f.Value == casted.Value
	} else {
		return false
	}
}

func (f *BinaryExpression) String() string {
	return fmt.Sprintf("%v %v %v", f.Field, f.Op, f.Value)
}

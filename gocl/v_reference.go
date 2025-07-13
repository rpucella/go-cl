package gocl

import (
	"fmt"
)

type vReference struct {
	content Value
}

func NewReference(v Value) Value {
	return &vReference{v}
}

func (v *vReference) Display() string {
	return fmt.Sprintf("#<ref %s>", v.content.Display())
}

func (v *vReference) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vReference) Apply(args []Value) (Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments %d to ref update", len(args))
	}
	if len(args) == 1 {
		v.content = args[0]
		return &vVoid{}, nil
	}
	return v.content, nil
}

func (v *vReference) str() string {
	return fmt.Sprintf("VReference[%s]", v.content.str())
}

func (v *vReference) IsAtom() bool {
	return false // ?
}

func (v *vReference) IsSymbol() bool {
	return false
}

func (v *vReference) IsCons() bool {
	return false
}

func (v *vReference) IsEmpty() bool {
	return false
}

func (v *vReference) IsNumber() bool {
	return false
}

func (v *vReference) IsBool() bool {
	return false
}

func (v *vReference) IsString() bool {
	return false
}

func (v *vReference) IsFunction() bool {
	return false
}

func (v *vReference) IsTrue() bool {
	return false
}

func (v *vReference) IsVoid() bool {
	return false
}

func (v *vReference) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *vReference) Type() string {
	return "reference"
}

func (v *vReference) AsInteger() (int, bool) {
	return 0, false
}

func (v *vReference) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vReference) AsString() (string, bool) {
	return "", false
}

func (v *vReference) AsSymbol() (string, bool) {
	return "", false
}

func (v *vReference) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vReference) AsReference() (Value, func(Value), bool) {
	update := func(cv Value) {
		v.content = cv
	}
	return v.content, update, true
}

func (v *vReference) SetReference(val Value) bool {
	v.content = val
	return true
}

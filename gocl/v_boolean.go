package gocl

import (
	"fmt"
)

type vBoolean struct {
	val bool
}

func NewBoolean(v bool) Value {
	return &vBoolean{v}
}

func (v *vBoolean) Display() string {
	if v.val {
		return "true"
	} else {
		return "false"
	}
}

func (v *vBoolean) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vBoolean) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vBoolean) str() string {
	if v.val {
		return "VBoolean[true]"
	} else {
		return "VBoolean[false]"
	}
}

func (v *vBoolean) IsAtom() bool {
	return true
}

func (v *vBoolean) IsSymbol() bool {
	return false
}

func (v *vBoolean) IsCons() bool {
	return false
}

func (v *vBoolean) IsEmpty() bool {
	return false
}

func (v *vBoolean) IsNumber() bool {
	return false
}

func (v *vBoolean) IsBool() bool {
	return true
}

func (v *vBoolean) IsString() bool {
	return false
}

func (v *vBoolean) IsFunction() bool {
	return false
}

func (v *vBoolean) IsTrue() bool {
	return v.val
}

func (v *vBoolean) IsVoid() bool {
	return false
}

func (v *vBoolean) IsEqual(vv Value) bool {
	b, ok := vv.AsBoolean()
	return ok && v.val == b
}

func (v *vBoolean) Type() string {
	return "bool"
}

func (v *vBoolean) AsInteger() (int, bool) {
	return 0, false
}

func (v *vBoolean) AsBoolean() (bool, bool) {
	return v.val, true
}

func (v *vBoolean) AsString() (string, bool) {
	return "", false
}

func (v *vBoolean) AsSymbol() (string, bool) {
	return "", false
}

func (v *vBoolean) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vBoolean) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vBoolean) SetReference(Value) bool {
	return false
}

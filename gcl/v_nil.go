package gcl

import (
	"fmt"
)

type vNil struct {
}

func NewNil() Value {
	return &vNil{}
}

func (v *vNil) Display() string {
	// figure out if this Is the right thing?
	return "#nil"
}

func (v *vNil) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vNil) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vNil) str() string {
	return fmt.Sprintf("VNil")
}

func (v *vNil) IsAtom() bool {
	return false
}

func (v *vNil) IsSymbol() bool {
	return false
}

func (v *vNil) IsCons() bool {
	return false
}

func (v *vNil) IsEmpty() bool {
	return false
}

func (v *vNil) IsNumber() bool {
	return false
}

func (v *vNil) IsBool() bool {
	return false
}

func (v *vNil) IsString() bool {
	return false
}

func (v *vNil) IsFunction() bool {
	return false
}

func (v *vNil) IsTrue() bool {
	return false
}

func (v *vNil) IsNil() bool {
	return true
}

func (v *vNil) IsEqual(vv Value) bool {
	return vv.IsNil()
}

func (v *vNil) Type() string {
	return "nil"
}

func (v *vNil) AsInteger() (int, bool) {
	return 0, false
}

func (v *vNil) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vNil) AsString() (string, bool) {
	return "", false
}

func (v *vNil) AsSymbol() (string, bool) {
	return "", false
}

func (v *vNil) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vNil) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vNil) SetReference(Value) bool {
	return false
}

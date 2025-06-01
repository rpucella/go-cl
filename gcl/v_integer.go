package gcl

import (
	"fmt"
)

type vInteger struct {
	val int
}

func NewInteger(v int) Value {
	return &vInteger{v}
}

func (v *vInteger) Display() string {
	return fmt.Sprintf("%d", v.val)
}

func (v *vInteger) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vInteger) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vInteger) str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *vInteger) IsAtom() bool {
	return true
}

func (v *vInteger) IsSymbol() bool {
	return false
}

func (v *vInteger) IsCons() bool {
	return false
}

func (v *vInteger) IsEmpty() bool {
	return false
}

func (v *vInteger) IsNumber() bool {
	return true
}

func (v *vInteger) IsBool() bool {
	return false
}

func (v *vInteger) IsString() bool {
	return false
}

func (v *vInteger) IsFunction() bool {
	return false
}

func (v *vInteger) IsTrue() bool {
	return v.val != 0
}

func (v *vInteger) IsVoid() bool {
	return false
}

func (v *vInteger) IsEqual(vv Value) bool {
	num, ok := vv.AsInteger()
	return ok && v.val == num
}

func (v *vInteger) Type() string {
	return "int"
}

func (v *vInteger) AsInteger() (int, bool) {
	return v.val, true
}

func (v *vInteger) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vInteger) AsString() (string, bool) {
	return "", false
}

func (v *vInteger) AsSymbol() (string, bool) {
	return "", false
}

func (v *vInteger) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vInteger) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vInteger) SetReference(Value) bool {
	return false
}

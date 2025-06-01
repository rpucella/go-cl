package gcl

import (
	"fmt"
)

type vEmpty struct {
}

func NewEmpty() Value {
	return &vEmpty{}
}

func (v *vEmpty) Display() string {
	return "()"
}

func (v *vEmpty) displayCDR() string {
	return ")"
}

func (v *vEmpty) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vEmpty) str() string {
	return fmt.Sprintf("VEmpty")
}

func (v *vEmpty) IsAtom() bool {
	return false
}

func (v *vEmpty) IsSymbol() bool {
	return false
}

func (v *vEmpty) IsCons() bool {
	return false
}

func (v *vEmpty) IsEmpty() bool {
	return true
}

func (v *vEmpty) IsNumber() bool {
	return false
}

func (v *vEmpty) IsBool() bool {
	return false
}

func (v *vEmpty) IsString() bool {
	return false
}

func (v *vEmpty) IsFunction() bool {
	return false
}

func (v *vEmpty) IsTrue() bool {
	return false
}

func (v *vEmpty) IsVoid() bool {
	return false
}

func (v *vEmpty) IsEqual(vv Value) bool {
	return vv.IsEmpty()
}

func (v *vEmpty) Type() string {
	return "list"
}

func (v *vEmpty) AsInteger() (int, bool) {
	return 0, false
}

func (v *vEmpty) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vEmpty) AsString() (string, bool) {
	return "", false
}

func (v *vEmpty) AsSymbol() (string, bool) {
	return "", false
}

func (v *vEmpty) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vEmpty) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vEmpty) SetReference(Value) bool {
	return false
}

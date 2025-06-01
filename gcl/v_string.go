package gcl

import (
	"fmt"
)

type vString struct {
	val string
}

func NewString(v string) Value {
	return &vString{v}
}

func (v *vString) Display() string {
	return "\"" + v.val + "\""
}

func (v *vString) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vString) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vString) str() string {
	return fmt.Sprintf("VString[%s]", v.str())
}

func (v *vString) IsAtom() bool {
	return true
}

func (v *vString) IsSymbol() bool {
	return false
}

func (v *vString) IsCons() bool {
	return false
}

func (v *vString) IsEmpty() bool {
	return false
}

func (v *vString) IsNumber() bool {
	return false
}

func (v *vString) IsBool() bool {
	return false
}

func (v *vString) IsString() bool {
	return true
}

func (v *vString) IsFunction() bool {
	return false
}

func (v *vString) IsTrue() bool {
	return (v.val != "")
}

func (v *vString) IsNil() bool {
	return false
}

func (v *vString) IsEqual(vv Value) bool {
	str, ok := vv.AsString()
	return ok && v.val == str
}

func (v *vString) Type() string {
	return "string"
}

func (v *vString) AsInteger() (int, bool) {
	return 0, false
}

func (v *vString) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vString) AsString() (string, bool) {
	return v.val, true
}

func (v *vString) AsSymbol() (string, bool) {
	return "", false
}

func (v *vString) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vString) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vString) SetReference(Value) bool {
	return false
}

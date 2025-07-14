package gocl

import (
	"fmt"
)

type vFlag struct {
	name string
}

func NewFlag(name string) Value {
	return &vFlag{name}
}

func (v *vFlag) Display() string {
	return fmt.Sprintf("--%s", v.name)
}

func (v *vFlag) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vFlag) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vFlag) str() string {
	return fmt.Sprintf("VFlag[%s]", v.name)
}

func (v *vFlag) IsAtom() bool {
	return true
}

func (v *vFlag) IsSymbol() bool {
	return true
}

func (v *vFlag) IsCons() bool {
	return false
}

func (v *vFlag) IsEmpty() bool {
	return false
}

func (v *vFlag) IsNumber() bool {
	return false
}

func (v *vFlag) IsBool() bool {
	return false
}

func (v *vFlag) IsString() bool {
	return false
}

func (v *vFlag) IsFunction() bool {
	return false
}

func (v *vFlag) IsTrue() bool {
	return true
}

func (v *vFlag) IsVoid() bool {
	return false
}

func (v *vFlag) IsEqual(vv Value) bool {
	name, ok := vv.AsFlag()
	return ok && v.name == name
}

func (v *vFlag) Type() string {
	return "symbol"
}

func (v *vFlag) AsInteger() (int, bool) {
	return 0, false
}

func (v *vFlag) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vFlag) AsString() (string, bool) {
	return "", false
}

func (v *vFlag) AsSymbol() (string, bool) {
	return "", false
}

func (v *vFlag) AsFlag() (string, bool) {
	return v.name, true
}

func (v *vFlag) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vFlag) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vFlag) SetReference(Value) bool {
	return false
}

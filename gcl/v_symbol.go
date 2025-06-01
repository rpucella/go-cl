package gcl

import (
	"fmt"
)

type vSymbol struct {
	name string
}

func NewSymbol(name string) Value {
	return &vSymbol{name}
}

func (v *vSymbol) Display() string {
	return v.name
}

func (v *vSymbol) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vSymbol) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vSymbol) str() string {
	return fmt.Sprintf("VSymbol[%s]", v.name)
}

func (v *vSymbol) IsAtom() bool {
	return true
}

func (v *vSymbol) IsSymbol() bool {
	return true
}

func (v *vSymbol) IsCons() bool {
	return false
}

func (v *vSymbol) IsEmpty() bool {
	return false
}

func (v *vSymbol) IsNumber() bool {
	return false
}

func (v *vSymbol) IsBool() bool {
	return false
}

func (v *vSymbol) IsString() bool {
	return false
}

func (v *vSymbol) IsFunction() bool {
	return false
}

func (v *vSymbol) IsTrue() bool {
	return true
}

func (v *vSymbol) IsNil() bool {
	return false
}

func (v *vSymbol) IsEqual(vv Value) bool {
	name, ok := vv.AsSymbol()
	return ok && v.name == name
}

func (v *vSymbol) Type() string {
	return "symbol"
}

func (v *vSymbol) AsInteger() (int, bool) {
	return 0, false
}

func (v *vSymbol) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vSymbol) AsString() (string, bool) {
	return "", false
}

func (v *vSymbol) AsSymbol() (string, bool) {
	return v.name, true
}

func (v *vSymbol) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vSymbol) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vSymbol) SetReference(Value) bool {
	return false
}

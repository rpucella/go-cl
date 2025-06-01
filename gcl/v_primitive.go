package gcl

import (
	"fmt"
)

type vPrimitive struct {
	name      string
	primitive func([]Value) (Value, error)
}

func NewPrimitive(name string, prim func([]Value) (Value, error)) Value {
	return &vPrimitive{name, prim}
}

func (v *vPrimitive) Display() string {
	return fmt.Sprintf("#<prim %s>", v.name)
}

func (v *vPrimitive) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vPrimitive) Apply(args []Value) (Value, error) {
	return v.primitive(args)
}

func (v *vPrimitive) str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *vPrimitive) IsAtom() bool {
	return false
}

func (v *vPrimitive) IsSymbol() bool {
	return false
}

func (v *vPrimitive) IsCons() bool {
	return false
}

func (v *vPrimitive) IsEmpty() bool {
	return false
}

func (v *vPrimitive) IsNumber() bool {
	return false
}

func (v *vPrimitive) IsBool() bool {
	return false
}

func (v *vPrimitive) IsString() bool {
	return false
}

func (v *vPrimitive) IsFunction() bool {
	return true
}

func (v *vPrimitive) IsTrue() bool {
	return true
}

func (v *vPrimitive) IsVoid() bool {
	return false
}

func (v *vPrimitive) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *vPrimitive) Type() string {
	return "fun"
}

func (v *vPrimitive) AsInteger() (int, bool) {
	return 0, false
}

func (v *vPrimitive) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vPrimitive) AsString() (string, bool) {
	return "", false
}

func (v *vPrimitive) AsSymbol() (string, bool) {
	return "", false
}

func (v *vPrimitive) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vPrimitive) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vPrimitive) SetReference(Value) bool {
	return false
}

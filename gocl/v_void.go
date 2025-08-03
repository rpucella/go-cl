package gocl

import (
	"fmt"
	"sync"
)

type vVoid struct {
}

var voidSingleton Value
var voidOnce sync.Once

func NewVoid() Value {
	voidOnce.Do(func() {
		voidSingleton = &vVoid{}
	})
	return voidSingleton
}

func (v *vVoid) Display() string {
	return "void"
}

func (v *vVoid) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vVoid) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vVoid) str() string {
	return fmt.Sprintf("VNil")
}

func (v *vVoid) IsAtom() bool {
	return false
}

func (v *vVoid) IsSymbol() bool {
	return false
}

func (v *vVoid) IsCons() bool {
	return false
}

func (v *vVoid) IsEmpty() bool {
	return false
}

func (v *vVoid) IsNumber() bool {
	return false
}

func (v *vVoid) IsBool() bool {
	return false
}

func (v *vVoid) IsString() bool {
	return false
}

func (v *vVoid) IsFunction() bool {
	return false
}

func (v *vVoid) IsTrue() bool {
	return false
}

func (v *vVoid) IsVoid() bool {
	return true
}

func (v *vVoid) IsEqual(vv Value) bool {
	return vv.IsVoid()
}

func (v *vVoid) Type() string {
	return "void"
}

func (v *vVoid) AsInteger() (int, bool) {
	return 0, false
}

func (v *vVoid) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vVoid) AsString() (string, bool) {
	return "", false
}

func (v *vVoid) AsSymbol() (string, bool) {
	return "", false
}

func (v *vVoid) AsFlag() (string, bool) {
	return "", false
}

func (v *vVoid) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vVoid) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vVoid) SetReference(Value) bool {
	return false
}

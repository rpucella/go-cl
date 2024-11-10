package glisp

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

func (v *vInteger) DisplayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vInteger) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vInteger) str() string {
	return fmt.Sprintf("VInteger[%d]", v.val)
}

func (v *vInteger) isAtom() bool {
	return true
}

func (v *vInteger) isSymbol() bool {
	return false
}

func (v *vInteger) isCons() bool {
	return false
}

func (v *vInteger) isEmpty() bool {
	return false
}

func (v *vInteger) isNumber() bool {
	return true
}

func (v *vInteger) isBool() bool {
	return false
}

func (v *vInteger) isString() bool {
	return false
}

func (v *vInteger) isFunction() bool {
	return false
}

func (v *vInteger) isTrue() bool {
	return v.val != 0
}

func (v *vInteger) isNil() bool {
	return false
}

func (v *vInteger) isEqual(vv Value) bool {
	num, ok := vv.asInteger()
	return ok && v.val == num
}

func (v *vInteger) typ() string {
	return "int"
}

func (v *vInteger) asInteger() (int, bool) {
	return v.val, true
}

func (v *vInteger) asBoolean() (bool, bool) {
	return false, false
}

func (v *vInteger) asString() (string, bool) {
	return "", false
}

func (v *vInteger) asSymbol() (string, bool) {
	return "", false
}

func (v *vInteger) asCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vInteger) asReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vInteger) setReference(Value) bool {
	return false
}

func (v *vInteger) asArray() ([]Value, bool) {
	return nil, false
}

func (v *vInteger) asDict() (map[string]Value, bool) {
	return nil, false
}

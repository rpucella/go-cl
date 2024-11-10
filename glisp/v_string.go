package glisp

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

func (v *vString) DisplayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vString) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vString) str() string {
	return fmt.Sprintf("VString[%s]", v.str())
}

func (v *vString) isAtom() bool {
	return true
}

func (v *vString) isSymbol() bool {
	return false
}

func (v *vString) isCons() bool {
	return false
}

func (v *vString) isEmpty() bool {
	return false
}

func (v *vString) isNumber() bool {
	return false
}

func (v *vString) isBool() bool {
	return false
}

func (v *vString) isString() bool {
	return true
}

func (v *vString) isFunction() bool {
	return false
}

func (v *vString) isTrue() bool {
	return (v.val != "")
}

func (v *vString) isNil() bool {
	return false
}

func (v *vString) isEqual(vv Value) bool {
	str, ok := vv.asString()
	return ok && v.val == str
}

func (v *vString) typ() string {
	return "string"
}

func (v *vString) asInteger() (int, bool) {
	return 0, false
}

func (v *vString) asBoolean() (bool, bool) {
	return false, false
}

func (v *vString) asString() (string, bool) {
	return v.val, true
}

func (v *vString) asSymbol() (string, bool) {
	return "", false
}

func (v *vString) asCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vString) asReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vString) setReference(Value) bool {
	return false
}

func (v *vString) asArray() ([]Value, bool) {
	return nil, false
}

func (v *vString) asDict() (map[string]Value, bool) {
	return nil, false
}

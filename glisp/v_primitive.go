package glisp

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

func (v *vPrimitive) DisplayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vPrimitive) apply(args []Value) (Value, error) {
	return v.primitive(args)
}

func (v *vPrimitive) str() string {
	return fmt.Sprintf("VPrimitive[%s]", v.name)
}

func (v *vPrimitive) isAtom() bool {
	return false
}

func (v *vPrimitive) isSymbol() bool {
	return false
}

func (v *vPrimitive) isCons() bool {
	return false
}

func (v *vPrimitive) isEmpty() bool {
	return false
}

func (v *vPrimitive) isNumber() bool {
	return false
}

func (v *vPrimitive) isBool() bool {
	return false
}

func (v *vPrimitive) isString() bool {
	return false
}

func (v *vPrimitive) isFunction() bool {
	return true
}

func (v *vPrimitive) isTrue() bool {
	return true
}

func (v *vPrimitive) isNil() bool {
	return false
}

func (v *vPrimitive) isEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *vPrimitive) typ() string {
	return "fun"
}

func (v *vPrimitive) asInteger() (int, bool) {
	return 0, false
}

func (v *vPrimitive) asBoolean() (bool, bool) {
	return false, false
}

func (v *vPrimitive) asString() (string, bool) {
	return "", false
}

func (v *vPrimitive) asSymbol() (string, bool) {
	return "", false
}

func (v *vPrimitive) asCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vPrimitive) asReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vPrimitive) setReference(Value) bool {
	return false
}

func (v *vPrimitive) asArray() ([]Value, bool) {
	return nil, false
}

func (v *vPrimitive) asDict() (map[string]Value, bool) {
	return nil, false
}

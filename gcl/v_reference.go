package gcl

import (
	"fmt"
)

type vReference struct {
	content Value
}

func NewReference(v Value) Value {
	return &vReference{v}
}

func (v *vReference) Display() string {
	return fmt.Sprintf("#<ref %s>", v.content.Display())
}

func (v *vReference) DisplayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vReference) apply(args []Value) (Value, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments %d to ref update", len(args))
	}
	if len(args) == 1 {
		v.content = args[0]
		return &vNil{}, nil
	}
	return v.content, nil
}

func (v *vReference) str() string {
	return fmt.Sprintf("VReference[%s]", v.content.str())
}

func (v *vReference) isAtom() bool {
	return false // ?
}

func (v *vReference) isSymbol() bool {
	return false
}

func (v *vReference) isCons() bool {
	return false
}

func (v *vReference) isEmpty() bool {
	return false
}

func (v *vReference) isNumber() bool {
	return false
}

func (v *vReference) isBool() bool {
	return false
}

func (v *vReference) isString() bool {
	return false
}

func (v *vReference) isFunction() bool {
	return false
}

func (v *vReference) isTrue() bool {
	return false
}

func (v *vReference) isNil() bool {
	return false
}

func (v *vReference) isEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *vReference) typ() string {
	return "reference"
}

func (v *vReference) asInteger() (int, bool) {
	return 0, false
}

func (v *vReference) asBoolean() (bool, bool) {
	return false, false
}

func (v *vReference) asString() (string, bool) {
	return "", false
}

func (v *vReference) asSymbol() (string, bool) {
	return "", false
}

func (v *vReference) asCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vReference) asReference() (Value, func(Value), bool) {
	update := func(cv Value) {
		v.content = cv
	}
	return v.content, update, true
}

func (v *vReference) setReference(val Value) bool {
	v.content = val
	return true
}

func (v *vReference) asArray() ([]Value, bool) {
	return nil, false
}

func (v *vReference) asDict() (map[string]Value, bool) {
	return nil, false
}

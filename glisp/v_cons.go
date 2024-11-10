package glisp

import (
	"fmt"
)

type vCons struct {
	head   Value
	tail   Value
	length int // Doesn't appear used.
}

func NewCons(head Value, tail Value) Value {
	return &vCons{head: head, tail: tail}
}

type MutableCons = *vCons

func NewMutableCons(head Value, tail Value) MutableCons {
	return &vCons{head: head, tail: tail}
}

func (v MutableCons) setTail(tail Value) {
	v.tail = tail
}

func (v *vCons) Display() string {
	return "(" + v.head.Display() + v.tail.DisplayCDR()
}

func (v *vCons) DisplayCDR() string {
	return " " + v.head.Display() + v.tail.DisplayCDR()
}

func (v *vCons) apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vCons) str() string {
	return fmt.Sprintf("VCons[%s %s]", v.head.str(), v.tail.str())
}

func (v *vCons) isAtom() bool {
	return false
}

func (v *vCons) isSymbol() bool {
	return false
}

func (v *vCons) isCons() bool {
	return true
}

func (v *vCons) isEmpty() bool {
	return false
}

func (v *vCons) isNumber() bool {
	return false
}

func (v *vCons) isBool() bool {
	return false
}

func (v *vCons) isString() bool {
	return false
}

func (v *vCons) isFunction() bool {
	return false
}

func (v *vCons) isTrue() bool {
	return true
}

func (v *vCons) isNil() bool {
	return false
}

func (v *vCons) isEqual(vv Value) bool {
	if _, _, ok := vv.asCons(); !ok {
		return false
	}
	var curr1 Value = v
	var curr2 Value = vv
	for head1, tail1, ok := v.asCons(); ok; head1, tail1, ok = tail1.asCons() {
		head2, tail2, ok := curr2.asCons()
		if !ok {
			return false
		}
		if !head1.isEqual(head2) {
			return false
		}
		curr1 = tail1
		curr2 = tail2
	}
	return curr1.isEqual(curr2) // should both be empty at the end
}

func (v *vCons) typ() string {
	return "list"
}

func (v *vCons) asInteger() (int, bool) {
	return 0, false
}

func (v *vCons) asBoolean() (bool, bool) {
	return false, false
}

func (v *vCons) asString() (string, bool) {
	return "", false
}

func (v *vCons) asSymbol() (string, bool) {
	return "", false
}

func (v *vCons) asCons() (Value, Value, bool) {
	return v.head, v.tail, true
}

func (v *vCons) asReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vCons) setReference(Value) bool {
	return false
}

func (v *vCons) asArray() ([]Value, bool) {
	return nil, false
}

func (v *vCons) asDict() (map[string]Value, bool) {
	return nil, false
}

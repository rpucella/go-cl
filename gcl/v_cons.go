package gcl

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
	return "(" + v.head.Display() + v.tail.displayCDR()
}

func (v *vCons) displayCDR() string {
	return " " + v.head.Display() + v.tail.displayCDR()
}

func (v *vCons) Apply(args []Value) (Value, error) {
	return nil, fmt.Errorf("Value %s not applicable", v.str())
}

func (v *vCons) str() string {
	return fmt.Sprintf("VCons[%s %s]", v.head.str(), v.tail.str())
}

func (v *vCons) IsAtom() bool {
	return false
}

func (v *vCons) IsSymbol() bool {
	return false
}

func (v *vCons) IsCons() bool {
	return true
}

func (v *vCons) IsEmpty() bool {
	return false
}

func (v *vCons) IsNumber() bool {
	return false
}

func (v *vCons) IsBool() bool {
	return false
}

func (v *vCons) IsString() bool {
	return false
}

func (v *vCons) IsFunction() bool {
	return false
}

func (v *vCons) IsTrue() bool {
	return true
}

func (v *vCons) IsVoid() bool {
	return false
}

func (v *vCons) IsEqual(vv Value) bool {
	if _, _, ok := vv.AsCons(); !ok {
		return false
	}
	var curr1 Value = v
	var curr2 Value = vv
	for head1, tail1, ok := v.AsCons(); ok; head1, tail1, ok = tail1.AsCons() {
		head2, tail2, ok := curr2.AsCons()
		if !ok {
			return false
		}
		if !head1.IsEqual(head2) {
			return false
		}
		curr1 = tail1
		curr2 = tail2
	}
	return curr1.IsEqual(curr2) // should both be empty at the end
}

func (v *vCons) Type() string {
	return "list"
}

func (v *vCons) AsInteger() (int, bool) {
	return 0, false
}

func (v *vCons) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vCons) AsString() (string, bool) {
	return "", false
}

func (v *vCons) AsSymbol() (string, bool) {
	return "", false
}

func (v *vCons) AsCons() (Value, Value, bool) {
	return v.head, v.tail, true
}

func (v *vCons) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vCons) SetReference(Value) bool {
	return false
}


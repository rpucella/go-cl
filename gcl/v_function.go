package gcl

import (
	"fmt"
	"strings"
)

type vFunction struct {
	params []string
	body   ast
	env    *Env
}

func NewFunction(params []string, body ast, env *Env) Value {
	return &vFunction{params, body, env}
}

func (v *vFunction) Display() string {
	return fmt.Sprintf("#<fun %s ...>", strings.Join(v.params, " "))
}

func (v *vFunction) displayCDR() string {
	panic(fmt.Sprintf("unchecked access to %s", v.str()))
}

func (v *vFunction) Apply(args []Value) (Value, error) {
	if len(v.params) != len(args) {
		return nil, fmt.Errorf("Wrong number of arguments to application to %s", v.str())
	}
	newEnv := layer(v.env, v.params, args)
	return v.body.eval(newEnv)
}

func (v *vFunction) str() string {
	return fmt.Sprintf("VFunction[[%s] %s]", strings.Join(v.params, " "), v.body.str())
}

func (v *vFunction) IsAtom() bool {
	return false
}

func (v *vFunction) IsSymbol() bool {
	return false
}

func (v *vFunction) IsCons() bool {
	return false
}

func (v *vFunction) IsEmpty() bool {
	return false
}

func (v *vFunction) IsNumber() bool {
	return false
}

func (v *vFunction) IsBool() bool {
	return false
}

func (v *vFunction) IsString() bool {
	return false
}

func (v *vFunction) IsFunction() bool {
	return true
}

func (v *vFunction) IsTrue() bool {
	return true
}

func (v *vFunction) IsVoid() bool {
	return false
}

func (v *vFunction) IsEqual(vv Value) bool {
	return v == vv // pointer equality
}

func (v *vFunction) Type() string {
	return "fun"
}

func (v *vFunction) AsInteger() (int, bool) {
	return 0, false
}

func (v *vFunction) AsBoolean() (bool, bool) {
	return false, false
}

func (v *vFunction) AsString() (string, bool) {
	return "", false
}

func (v *vFunction) AsSymbol() (string, bool) {
	return "", false
}

func (v *vFunction) AsCons() (Value, Value, bool) {
	return nil, nil, false
}

func (v *vFunction) AsReference() (Value, func(Value), bool) {
	return nil, nil, false
}

func (v *vFunction) SetReference(Value) bool {
	return false
}

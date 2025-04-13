package gcl

type Value interface {
	Display() string
	DisplayCDR() string

	asInteger() (int, bool)
	asBoolean() (bool, bool)
	asString() (string, bool)
	asSymbol() (string, bool)
	asCons() (Value, Value, bool)
	asReference() (Value, func(Value), bool)
	setReference(Value) bool
	asArray() ([]Value, bool)
	asDict() (map[string]Value, bool)

	apply([]Value) (Value, error)
	str() string
	isAtom() bool
	isEmpty() bool
	isTrue() bool
	isNil() bool
	isFunction() bool
	//isEq() bool    -- don't think we need pointer equality for now - = is enough?
	isEqual(Value) bool
	typ() string
}

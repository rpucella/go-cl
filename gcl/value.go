package gcl

type Value interface {
	Display() string

	AsInteger() (int, bool)
	AsBoolean() (bool, bool)
	AsString() (string, bool)
	AsSymbol() (string, bool)
	AsCons() (Value, Value, bool)
	AsReference() (Value, func(Value), bool)
	SetReference(Value) bool

	Apply([]Value) (Value, error)
	IsAtom() bool
	IsEmpty() bool
	IsTrue() bool
	IsVoid() bool
	IsFunction() bool
	//isEq() bool    -- don't think we need pointer equality for now - = is enough?
	IsEqual(Value) bool
	Type() string

	// Internal operations.
	displayCDR() string
	str() string
}

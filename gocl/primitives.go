package gocl

import "fmt"
import "strings"

type Primitive struct {
	name string
	min  int
	max  int // <0 for no max #
	prim func(string, []Value) (Value, error)
}

type Command struct {
	// A command is like a primitive that also takes optional flags as inputs.
	name string
	min int
	max int  // <0 for no max #
	prim func(string, []Value, map[string]bool)
}

func listLength(v Value) int {
	result := 0
	for _, next, ok := v.AsCons(); ok; _, next, ok = next.AsCons() {
		result += 1
	}
	return result
}

func listAppend(v1 Value, v2 Value) Value {
	var result Value = nil
	var current_result MutableCons = nil
	for head, next, ok := v1.AsCons(); ok; head, next, ok = next.AsCons() {
		cell := NewMutableCons(head, nil)
		if current_result == nil {
			result = cell
		} else {
			current_result.setTail(cell)
		}
		current_result = cell
	}
	if current_result == nil {
		return v2
	}
	current_result.setTail(v2)
	return result
}

func allConses(vs []Value) bool {
	for _, v := range vs {
		if _, _, ok := v.AsCons(); !ok {
			return false
		}
	}
	return true
}

func corePrimitives() map[string]Value {
	bindings := map[string]Value{}
	for _, d := range primList {
		bindings[d.name] = NewPrimitive(d.name, MakePrimitive(d))
	}
	return bindings
}

// func MakeCommand(d Primitive)

func MakePrimitive(d Primitive) func([]Value) (Value, error) {
	f := func(args []Value) (Value, error) {
		if err := checkMinArgs(d.name, args, d.min); err != nil {
			return nil, err
		}
		if d.max >= 0 {
			if err := checkMaxArgs(d.name, args, d.max); err != nil {
				return nil, err
			}
		}
		return d.prim(d.name, args)
	}
	return f
}

func checkArgType(name string, arg Value, pred func(Value) bool) error {
	if !pred(arg) {
		return fmt.Errorf("%s - wrong argument type %s", name, arg.Type())
	}
	return nil
}

func checkArgTypeB(name string, arg Value, ok bool) error {
	if !ok {
		return fmt.Errorf("%s - wrong argument type %s", name, arg.Type())
	}
	return nil
}

func checkMinArgs(name string, args []Value, n int) error {
	if len(args) < n {
		return fmt.Errorf("%s - too few arguments %d", name, len(args))
	}
	return nil
}

func checkMaxArgs(name string, args []Value, n int) error {
	if len(args) > n {
		return fmt.Errorf("%s - too many arguments %d", name, len(args))
	}
	return nil
}

func checkExactArgs(name string, args []Value, n int) error {
	if len(args) != n {
		return fmt.Errorf("%s - wrong number of arguments %d", name, len(args))
	}
	return nil
}

func isInt(v Value) bool {
	_, ok := v.AsInteger()
	return ok
}

func isString(v Value) bool {
	_, ok := v.AsString()
	return ok
}

func isFunction(v Value) bool {
	return v.IsFunction()
}

func isList(v Value) bool {
	if _, _, ok := v.AsCons(); ok {
		return true
	}
	return v.IsEmpty()
}

func mkNumPredicate(pred func(int, int) bool) func(string, []Value) (Value, error) {
	return func(name string, args []Value) (Value, error) {
		if err := checkExactArgs(name, args, 2); err != nil {
			return nil, err
		}
		i1, ok := args[0].AsInteger()
		if err := checkArgTypeB(name, args[0], ok); err != nil {
			return nil, err
		}
		i2, ok := args[1].AsInteger()
		if err := checkArgTypeB(name, args[1], ok); err != nil {
			return nil, err
		}
		return NewBoolean(pred(i1, i2)), nil
	}
}

var primList = []Primitive{

	Primitive{
		"type", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewSymbol(args[0].Type()), nil
		},
	},

	Primitive{
		"+", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := 0
			for _, arg := range args {
				i1, ok := arg.AsInteger()
				if err := checkArgTypeB(name, arg, ok); err != nil {
					return nil, err
				}
				v += i1
			}
			return NewInteger(v), nil
		},
	},

	Primitive{
		"*", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := 1
			for _, arg := range args {
				i1, ok := arg.AsInteger()
				if err := checkArgTypeB(name, arg, ok); err != nil {
					return nil, err
				}
				v *= i1
			}
			return NewInteger(v), nil
		},
	},

	Primitive{
		"-", 1, -1,
		func(name string, args []Value) (Value, error) {
			v, ok := args[0].AsInteger()
			if err := checkArgTypeB(name, args[0], ok); err != nil {
				return nil, err
			}
			if len(args) > 1 {
				for _, arg := range args[1:] {
					i1, ok := arg.AsInteger()
					if err := checkArgTypeB(name, arg, ok); err != nil {
						return nil, err
					}
					v -= i1
				}
			} else {
				v = -v
			}
			return NewInteger(v), nil
		},
	},

	Primitive{"=", 2, -1,
		func(name string, args []Value) (Value, error) {
			var reference Value = args[0]
			for _, v := range args[1:] {
				if !reference.IsEqual(v) {
					return NewBoolean(false), nil
				}
			}
			return NewBoolean(true), nil
		},
	},

	Primitive{"<", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 < n2 }),
	},

	Primitive{"<=", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 <= n2 }),
	},

	Primitive{">", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 > n2 }),
	},

	Primitive{">=", 2, 2,
		mkNumPredicate(func(n1 int, n2 int) bool { return n1 >= n2 }),
	},

	Primitive{"not", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewBoolean(!args[0].IsTrue()), nil
		},
	},

	Primitive{
		"string-append", 0, -1,
		func(name string, args []Value) (Value, error) {
			v := ""
			for _, arg := range args {
				str, ok := arg.AsString()
				if err := checkArgTypeB(name, arg, ok); err != nil {
					return nil, err
				}
				v += str
			}
			return NewString(v), nil
		},
	},

	Primitive{"string-length", 1, 1,
		func(name string, args []Value) (Value, error) {
			str, ok := args[0].AsString()
			if err := checkArgTypeB(name, args[0], ok); err != nil {
				return nil, err
			}
			return NewInteger(len(str)), nil
		},
	},

	Primitive{"string-lower", 1, 1,
		func(name string, args []Value) (Value, error) {
			str, ok := args[0].AsString()
			if err := checkArgTypeB(name, args[0], ok); err != nil {
				return nil, err
			}
			return NewString(strings.ToLower(str)), nil
		},
	},

	Primitive{"string-upper", 1, 1,
		func(name string, args []Value) (Value, error) {
			str, ok := args[0].AsString()
			if err := checkArgTypeB(name, args[0], ok); err != nil {
				return nil, err
			}
			return NewString(strings.ToUpper(str)), nil
		},
	},

	Primitive{"string-substring", 1, 3,
		func(name string, args []Value) (Value, error) {
			str, ok := args[0].AsString()
			if err := checkArgTypeB(name, args[0], ok); err != nil {
				return nil, err
			}
			start := 0
			end := len(str)
			if len(args) > 2 {
				i1, ok := args[2].AsInteger()
				if err := checkArgTypeB(name, args[2], ok); err != nil {
					return nil, err
				}
				end = min(i1, end)
			}
			if len(args) > 1 {
				i1, ok := args[1].AsInteger()
				if err := checkArgTypeB(name, args[1], ok); err != nil {
					return nil, err
				}
				start = max(i1, start)
			}
			// or perhaps raise an exception
			if end < start {
				return NewString(""), nil
			}
			return NewString(str[start:end]), nil
		},
	},

	Primitive{"apply", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			arguments := make([]Value, listLength(args[1]))
			current := args[1]
			for i := range arguments {
				head, tail, _ := current.AsCons() // isList before checked ok.
				arguments[i] = head
				current = tail
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return args[0].Apply(arguments)
		},
	},

	Primitive{"cons", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			return NewCons(args[0], args[1]), nil
		},
	},

	Primitive{
		"append", 0, -1,
		func(name string, args []Value) (Value, error) {
			if len(args) == 0 {
				return NewEmpty(), nil
			}
			if err := checkArgType(name, args[len(args)-1], isList); err != nil {
				return nil, err
			}
			result := args[len(args)-1]
			for i := len(args) - 2; i >= 0; i -= 1 {
				if err := checkArgType(name, args[i], isList); err != nil {
					return nil, err
				}
				result = listAppend(args[i], result)
			}
			return result, nil
		},
	},

	Primitive{"reverse", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			var result Value = NewEmpty()
			current := args[0]
			for head, next, ok := args[0].AsCons(); ok; head, next, ok = next.AsCons() {
				result = NewCons(head, result)
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	Primitive{"head", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if args[0].IsEmpty() {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			head, _, _ := args[0].AsCons()
			return head, nil
		},
	},

	Primitive{"tail", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			if args[0].IsEmpty() {
				return nil, fmt.Errorf("%s - empty list argument", name)
			}
			_, tail, _ := args[0].AsCons()
			return tail, nil
		},
	},

	Primitive{"list", 0, -1,
		func(name string, args []Value) (Value, error) {
			var result Value = NewEmpty()
			for i := len(args) - 1; i >= 0; i -= 1 {
				result = NewCons(args[i], result)
			}
			return result, nil
		},
	},

	Primitive{"length", 1, 1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			count := 0
			current := args[0]
			for _, next, ok := args[0].AsCons(); ok; _, next, ok = next.AsCons() {
				count += 1
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return NewInteger(count), nil
		},
	},

	Primitive{"nth", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isList); err != nil {
				return nil, err
			}
			idx, ok := args[1].AsInteger()
			if err := checkArgTypeB(name, args[1], ok); err != nil {
				return nil, err
			}
			if idx >= 0 {
				for head, next, ok := args[0].AsCons(); ok; head, next, ok = next.AsCons() {
					if idx == 0 {
						return head, nil
					} else {
						idx -= 1
					}
				}
			}
			return nil, fmt.Errorf("%s - index %d out of bound", name, idx)
		},
	},

	Primitive{"map", 2, -1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			var result Value = nil
			var current_result MutableCons = nil
			currents := make([]Value, len(args)-1)
			firsts := make([]Value, len(args)-1)
			for i := range args[1:] {
				currents[i] = args[i+1]
			}
			for allConses(currents) {
				for i := range currents {
					head, _, _ := currents[i].AsCons()
					firsts[i] = head
				}
				v, err := args[0].Apply(firsts)
				if err != nil {
					return nil, err
				}
				cell := NewMutableCons(v, nil)
				if current_result == nil {
					result = cell
				} else {
					current_result.setTail(cell)
				}
				current_result = cell
				for i := range currents {
					_, tail, _ := currents[i].AsCons()
					currents[i] = tail
				}
			}
			if current_result == nil {
				return NewEmpty(), nil
			}
			current_result.setTail(NewEmpty())
			return result, nil
		},
	},

	Primitive{"for", 2, -1,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			// TODO - allow different types in the same iteration!
			for i := range args[1:] {
				if err := checkArgType(name, args[i+1], isList); err != nil {
					return nil, err
				}
			}
			currents := make([]Value, len(args)-1)
			firsts := make([]Value, len(args)-1)
			for i := range args[1:] {
				currents[i] = args[i+1]
			}
			for allConses(currents) {
				for i := range currents {
					head, _, _ := currents[i].AsCons()
					firsts[i] = head
				}
				_, err := args[0].Apply(firsts)
				if err != nil {
					return nil, err
				}
				for i := range currents {
					_, tail, _ := currents[i].AsCons()
					currents[i] = tail
				}
			}
			return NewVoid(), nil
		},
	},

	Primitive{"filter", 2, 2,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var result Value = nil
			var current_result MutableCons = nil
			current := args[1]
			for head, next, ok := args[1].AsCons(); ok; head, next, ok = next.AsCons() {
				v, err := args[0].Apply([]Value{head})
				if err != nil {
					return nil, err
				}
				if v.IsTrue() {
					cell := NewMutableCons(head, nil)
					if current_result == nil {
						result = cell
					} else {
						current_result.setTail(cell)
					}
					current_result = cell
				}
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			if current_result == nil {
				return NewEmpty(), nil
			}
			current_result.setTail(NewEmpty())
			return result, nil
		},
	},

	Primitive{"foldr", 3, 3,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			var temp Value = NewEmpty()
			// first reverse the list
			current := args[1]
			for head, next, ok := args[1].AsCons(); ok; head, next, ok = next.AsCons() {
				temp = NewCons(head, temp)
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			// then fold it
			result := args[2]
			current = temp
			for head, next, ok := temp.AsCons(); ok; head, next, ok = next.AsCons() {
				v, err := args[0].Apply([]Value{head, result})
				if err != nil {
					return nil, err
				}
				result = v
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	Primitive{"foldl", 3, 3,
		func(name string, args []Value) (Value, error) {
			if err := checkArgType(name, args[0], isFunction); err != nil {
				return nil, err
			}
			if err := checkArgType(name, args[1], isList); err != nil {
				return nil, err
			}
			result := args[2]
			current := args[1]
			for head, next, ok := args[1].AsCons(); ok; head, next, ok = next.AsCons() {
				v, err := args[0].Apply([]Value{result, head})
				if err != nil {
					return nil, err
				}
				result = v
				current = next
			}
			if !current.IsEmpty() {
				return nil, fmt.Errorf("%s - malformed list", name)
			}
			return result, nil
		},
	},

	Primitive{"ref", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewReference(args[0]), nil
		},
	},

	Primitive{"empty?", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewBoolean(args[0].IsEmpty()), nil
		},
	},

	Primitive{"cons?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, _, ok := args[0].AsCons()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"list?", 1, 1,
		func(name string, args []Value) (Value, error) {
			result := false
			if _, _, ok := args[0].AsCons(); ok {
				result = true
			} else {
				result = args[0].IsEmpty()
			}
			return NewBoolean(result), nil
		},
	},

	Primitive{"number?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, ok := args[0].AsInteger()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"ref?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, _, ok := args[0].AsReference()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"boolean?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, ok := args[0].AsBoolean()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"string?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, ok := args[0].AsString()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"symbol?", 1, 1,
		func(name string, args []Value) (Value, error) {
			_, ok := args[0].AsSymbol()
			return NewBoolean(ok), nil
		},
	},

	Primitive{"function?", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewBoolean(args[0].IsFunction()), nil
		},
	},

	Primitive{"void?", 1, 1,
		func(name string, args []Value) (Value, error) {
			return NewBoolean(args[0].IsVoid()), nil
		},
	},
}

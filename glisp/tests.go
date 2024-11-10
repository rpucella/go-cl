package glisp

import "fmt"

func test() {
	test_value_10()
	test_value_plus()
	test_literal()
	test_lookup()
	test_apply()
	test_if()
	test_lists()
	test_read()
}

func getInt(v Value) int {
	i, ok := v.asInteger()
	if ok {
		return i
	}
	panic("not an integer")
}

func primitiveAdd(args []Value) (Value, error) {
	var result int
	for _, val := range args {
		result += getInt(val)
	}
	return &vInteger{result}, nil
}

func primitiveMult(args []Value) (Value, error) {
	var result int = 1
	for _, val := range args {
		result *= getInt(val)
	}
	return &vInteger{result}, nil
}

func sampleEnv() *Env {
	current := map[string]Value{
		"a": &vInteger{10},
		"b": &vInteger{20},
		"+": &vPrimitive{"+", primitiveAdd},
		"*": &vPrimitive{"*", primitiveMult},
		"t": &vBoolean{true},
		"f": &vBoolean{false},
	}
	env := &Env{bindings: current}
	return env
}

func test_value_10() {
	var v1 Value = &vInteger{10}
	fmt.Println(v1.str(), "->", getInt(v1))
}

func test_value_plus() {
	var v1 Value = &vInteger{10}
	var v2 Value = &vInteger{20}
	var v3 Value = &vInteger{30}
	var vp Value = &vPrimitive{"+", primitiveAdd}
	var args []Value = []Value{v1, v2, v3}
	vr, _ := vp.apply(args)
	fmt.Println(vp.str(), "->", getInt(vr))
}

func evalDisplay(e ast, env *Env) string {
	v, _ := e.eval(env)
	return v.Display()
}

func test_literal() {
	v1 := &vInteger{10}
	e1 := &astLiteral{v1}
	fmt.Println(e1.str(), "->", evalDisplay(e1, nil))
	v2 := &vBoolean{true}
	e2 := &astLiteral{v2}
	fmt.Println(e2.str(), "->", evalDisplay(e2, nil))
}

func test_lookup() {
	env := sampleEnv()
	e1 := &astId{"a"}
	fmt.Println(e1.str(), "->", evalDisplay(e1, env))
	e2 := &astId{"+"}
	fmt.Println(e2.str(), "->", evalDisplay(e2, env))
}

func test_apply() {
	env := sampleEnv()
	e1 := &astId{"a"}
	e2 := &astId{"b"}
	args := []ast{e1, e2}
	e3 := &astApply{&astId{"+"}, args}
	fmt.Println(e3.str(), "->", evalDisplay(e3, env))
}

func test_if() {
	env := sampleEnv()
	e1 := &astIf{&astId{"t"}, &astId{"a"}, &astId{"b"}}
	fmt.Println(e1.str(), "->", evalDisplay(e1, env))
	e2 := &astIf{&astId{"f"}, &astId{"a"}, &astId{"b"}}
	fmt.Println(e2.str(), "->", evalDisplay(e2, env))
}

func test_read() {
	v1, _, _ := read("123")
	fmt.Println(v1.str(), "->", v1.Display())
	v2, _, _ := read("a")
	fmt.Println(v2.str(), "->", v2.Display())
	v6, _, _ := read("+")
	fmt.Println(v6.str(), "->", v6.Display())
	v3, _, _ := read("(+ 33 a)")
	fmt.Println(v3.str(), "->", v3.Display())
	v4, _, _ := read("(+ 33 (+ a b))")
	fmt.Println(v4.str(), "->", v4.Display())
	v5, _, _ := read("(this is my life)")
	fmt.Println(v5.str(), "->", v5.Display())
}

func test_lists() {
	var v Value = &vEmpty{}
	v = &vCons{head: &vInteger{33}, tail: v}
	v = &vCons{head: &vInteger{66}, tail: v}
	v = &vCons{head: &vInteger{99}, tail: v}
	fmt.Println(v.str(), "->", v.Display())
}

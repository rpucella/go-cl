package gcl

import (
	"fmt"
)

type Engine struct {
	env *Env
}

func NewEngine() Engine {
	coreBindings := corePrimitives()
	coreBindings["true"] = NewBoolean(true)
	coreBindings["false"] = NewBoolean(false)
	env := &Env{bindings: coreBindings, previous: nil}
	return Engine{env}
}


// TODO: engine.Read()
// TODO: engine.Eval()
// TODO: engine.ReadEval()

// TODO: engine.DefConstant()
// TODO: engine.DefFunction()
// TODO: engine.DefMacro()

// TODO: handle multiline values

// TODO: make prompt a function (of what?)

// TODO: what do we export? Engine, Value

// TODO: decide if we need an interface? when do we create interfaces?

// TODO: figure out how to create documentation
// TODO: exceptions? Simple catch / throw?
// (string x)
// (number y)
// (bool z)

// Put direct in the engine itself?

func (e Engine) Read(line string, prefix string) (Value, error) {
	text := prefix + " " + line
	// TODO: distinguish error from "incomplete"?
	v, _, err := read(text)
	return v, err
}

func (e Engine) ReadList(line string) (Value, error) {
	text := line
	v, _, err := readList(text)
	return v, err
}

func (e Engine) ProcessDeclaration(v Value) (string, error) {
	// TODO: IsDeclaration that also returns the declaration.
	// engine.Def(decl) operation
	d, err := parseDef(v)
	if err != nil {
		return "", err
	}
	if d != nil {
		if d.typ == DEF_FUNCTION {
			update(e.env, d.name, &vFunction{d.params, d.body, e.env})
			return d.name, nil
		}
		if d.typ == DEF_VALUE {
			v, err := d.body.eval(e.env)
			if err != nil {
				fmt.Println("EVAL ERROR -", err.Error())
			}
			update(e.env, d.name, v)
			return d.name, nil
		}
		return "", fmt.Errorf("DECLARE ERROR - unknow declaration type %s", d.typ)
	}
	return "", nil
}

func (e Engine) AddPrimitive(name string, min int, max int, p func(string, []Value) (Value, error)) {
	pr := NewPrimitive(name, MakePrimitive(Primitive{name, min, max, p}))
	update(e.env, name, pr)
}

func (e Engine) Eval(sexpr Value) (Value, error) {
	expr, err := parseExpr(sexpr)
	if err != nil {
		fmt.Println()
		return nil, fmt.Errorf("PARSE ERROR - %w", err.Error())
	}
	///fmt.Println("expr =", e.str())
	v, err := expr.eval(e.env)
	if err != nil {
		return nil, fmt.Errorf("EVAL ERROR - %w", err.Error())
	}
	return v, nil
}

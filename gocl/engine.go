package gocl

import (
	"fmt"
	"strings"
)

type Engine struct {
	env *Env
	commands []string
	helpArgs []string
}

func NewEngine() Engine {
	coreBindings := corePrimitives()
	// These are not special. They're just convenience.
	coreBindings["empty"] = NewEmpty()
	coreBindings["void"] = NewVoid()
	env := &Env{bindings: coreBindings, previous: nil}
	commands := make([]string, 0)
	helpArgs := make([]string, 0)
	return Engine{env, commands, helpArgs}
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

func (e *Engine) Read(line string, prefix string) (Value, error) {
	text := prefix + " " + line
	// TODO: distinguish error from "incomplete"?
	v, _, err := read(text)
	return v, err
}

func (e *Engine) ReadList(line string) (Value, error) {
	text := line
	v, _, err := readList(text)
	return v, err
}

func (e *Engine) ProcessDeclaration(v Value) (string, error) {
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
				return "", fmt.Errorf("EVAL ERROR - %w", err)
			}
			update(e.env, d.name, v)
			return d.name, nil
		}
		return "", fmt.Errorf("DECLARE ERROR - unknow declaration type %s", d.typ)
	}
	return "", nil
}

func (e *Engine) AddPrimitive(name string, min int, max int, p func(string, []Value) (Value, error)) {
	pr := NewPrimitive(name, MakePrimitive(Primitive{name, min, max, p}))
	update(e.env, name, pr)
}

func (e *Engine) AddBinding(name string, val Value) {
	update(e.env, name, val)
}

func (e *Engine) Eval(sexpr Value) (Value, error) {
	expr, err := parseExpr(sexpr)
	if err != nil {
		fmt.Println()
		return nil, fmt.Errorf("PARSE ERROR - %w", err)
	}
	///fmt.Println("expr =", e.str())
	v, err := expr.eval(e.env)
	if err != nil {
		return nil, fmt.Errorf("EVAL ERROR - %w", err)
	}
	return v, nil
}


func (e *Engine) AddCommand(name string, helpArgs string, min int, max int, flags []string, p func(string, []Value) (Value, error)) {
	pr := NewPrimitive(name, MakePrimitive(Primitive{name, min, max, p}))
	e.commands = append(e.commands, name)
	e.helpArgs = append(e.helpArgs, helpArgs)
	update(e.env, name, pr)
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func (e *Engine) ReadCommand(line string) (Value, error) {
	v, _, err := readList(line)
	if err != nil {
		return nil, err
	}
	head, tail, ok := v.AsCons()
	if !ok {
		return nil, fmt.Errorf("malformed command")
	}
	comm, ok := head.AsSymbol()
	if !ok || !contains(e.commands, comm) {
		return nil, fmt.Errorf("unknown command: %s", comm)
	}
	// Transform every symbol into a string!
	args := make([]Value, 0)
	curr := tail
	for !curr.IsEmpty() {
		x, y, ok := curr.AsCons()
		if !ok {
			return nil, fmt.Errorf("malformed command")
		}
		sym, ok := x.AsSymbol()
		if ok {
			args = append(args, NewString(sym))
		} else {
			if !x.IsAtom() {
				return nil, fmt.Errorf("malformed command")
			}
			args = append(args, x)
		}
		curr = y
	}
	result := NewEmpty()
	for i := len(args) - 1; i >= 0; i-- {
		result = NewCons(args[i], result)
	}
	return NewCons(head, result), nil
}

func (e *Engine) ReadCommandWords(words []string) (Value, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("no command given")
	}
	result := NewEmpty()
	for i := len(words) - 1; i >= 0; i-- {
		w := strings.TrimSpace(words[i])
		v, rest, err := read(strings.TrimSpace(w))
		if err != nil {
			return nil, err
		}
		if i == 0 {
			comm, ok := v.AsSymbol()
			if !ok || !contains(e.commands, comm) {
				return nil, fmt.Errorf("unknown command: %s", comm)
			}
			result = NewCons(v, result)
		} else if rest != "" {
			result = NewCons(NewString(w), result)
		} else if sym, ok := v.AsSymbol(); ok {
			result = NewCons(NewString(sym), result)
		} else if v.IsAtom() {
			result = NewCons(v, result)
		} else {
			return nil, fmt.Errorf("malformed command")
		}
	}
	return result, nil
}

func (e *Engine) AddDefaultHelpCommand() {
	prim := func(name string, args []Value) (Value, error) {
		fmt.Println("Available commands:")
		for i, comm := range e.commands {
			fmt.Printf("  %s %s\n", comm, e.helpArgs[i])
		}
		return NewVoid(), nil
	}
	e.AddCommand("help", "", 0, 0, nil, prim)
}

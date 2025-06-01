package gcl

import "fmt"

type Env struct {
	bindings map[string]Value
	previous *Env
}

func find(env *Env, name string) (Value, error) {
	current := env
	for current != nil {
		val, ok := current.bindings[name]
		if ok {
			return val, nil
		}
		current = current.previous
	}
	return nil, fmt.Errorf("no such identifier %s", name)
}

func update(env *Env, name string, v Value) {
	env.bindings[name] = v
}

func layer(env *Env, names []string, values []Value) *Env {
	// if values is nil or smaller than names, then
	// remaining names are bound to void
	bindings := map[string]Value{}
	for i, name := range names {
		if values != nil && i < len(values) {
			bindings[name] = values[i]
		} else {
			bindings[name] = &vVoid{}
		}
	}
	return &Env{bindings: bindings, previous: env}
}

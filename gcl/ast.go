package gcl

import "fmt"
import "errors"
import "strings"

const DEF_VALUE = 0
const DEF_FUNCTION = 1

type astDef struct {
	name   string
	typ    int
	params []string
	body   ast
}

type ast interface {
	eval(*Env) (Value, error)
	evalPartial(*Env) (*partialResult, error)
	str() string
}

type partialResult struct {
	exp ast
	env *Env
	val Value // val is null when the result is still partial
}

type astLiteral struct {
	val Value
}

type astId struct {
	name string
}

type astIf struct {
	cnd ast
	thn ast
	els ast
}

type astApply struct {
	fn   ast
	args []ast
}

type astQuote struct {
	val Value
}

type astLetRec struct {
	names  []string
	params [][]string
	bodies []ast
	body   ast
}

func defaultEvalPartial(e ast, env *Env) (*partialResult, error) {
	// Partial evaluation
	// Sometimes return an expression to evaluate next along
	// with an environment for evaluation.
	// val is null when the result is in fact a value.

	v, err := e.eval(env)
	if err != nil {
		return nil, err
	}
	return &partialResult{nil, nil, v}, nil
}

func defaultEval(e ast, env *Env) (Value, error) {
	// evaluation with tail call optimization
	var currExp ast = e
	currEnv := env
	for {
		partial, err := currExp.evalPartial(currEnv)
		if err != nil {
			return nil, err
		}
		if partial.val != nil {
			return partial.val, nil
		}
		currExp = partial.exp
		currEnv = partial.env
	}
}

func (e *astLiteral) eval(env *Env) (Value, error) {
	return e.val, nil
}

func (e *astLiteral) evalPartial(env *Env) (*partialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *astLiteral) str() string {
	return fmt.Sprintf("astLiteral[%s]", e.val.str())
}

func (e *astId) eval(env *Env) (Value, error) {
	return find(env, e.name)
}

func (e *astId) evalPartial(env *Env) (*partialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *astId) str() string {
	return fmt.Sprintf("astId[%s]", e.name)
}

func (e *astIf) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *astIf) evalPartial(env *Env) (*partialResult, error) {
	c, err := e.cnd.eval(env)
	if err != nil {
		return nil, err
	}
	if c.IsTrue() {
		return &partialResult{e.thn, env, nil}, nil
	} else {
		return &partialResult{e.els, env, nil}, nil
	}
}

func (e *astIf) str() string {
	return fmt.Sprintf("astIf[%s %s %s]", e.cnd.str(), e.thn.str(), e.els.str())
}

func (e *astApply) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *astApply) evalPartial(env *Env) (*partialResult, error) {
	f, err := e.fn.eval(env)
	if err != nil {
		return nil, err
	}
	args := make([]Value, len(e.args))
	for i := range args {
		args[i], err = e.args[i].eval(env)
		if err != nil {
			return nil, err
		}
	}
	if ff, ok := f.(*vFunction); ok {
		if len(ff.params) != len(args) {
			return nil, fmt.Errorf("Wrong number of arguments to application to %s", ff.str())
		}
		newEnv := layer(ff.env, ff.params, args)
		return &partialResult{ff.body, newEnv, nil}, nil
	}
	v, err := f.Apply(args)
	if err != nil {
		return nil, err
	}
	return &partialResult{nil, nil, v}, nil
}

func (e *astApply) str() string {
	strArgs := ""
	for _, item := range e.args {
		strArgs += " " + item.str()
	}
	return fmt.Sprintf("astApply[%s%s]", e.fn.str(), strArgs)
}

func (e *astQuote) eval(env *Env) (Value, error) {
	return e.val, nil
}

func (e *astQuote) evalPartial(env *Env) (*partialResult, error) {
	return defaultEvalPartial(e, env)
}

func (e *astQuote) str() string {
	return fmt.Sprintf("astQuote[%s]", e.val.str())
}

func (e *astLetRec) eval(env *Env) (Value, error) {
	return defaultEval(e, env)
}

func (e *astLetRec) evalPartial(env *Env) (*partialResult, error) {
	if len(e.names) != len(e.params) || len(e.names) != len(e.bodies) {
		return nil, errors.New("malformed letrec (names, params, bodies)")
	}
	// create the environment that we'll share across the definitions
	// all names initially allocated #nil
	newEnv := layer(env, e.names, nil)
	for i, name := range e.names {
		update(newEnv, name, &vFunction{e.params[i], e.bodies[i], newEnv})
	}
	return &partialResult{e.body, newEnv, nil}, nil
}

func (e *astLetRec) str() string {
	bindings := make([]string, len(e.names))
	for i := range e.names {
		params := strings.Join(e.params[i], " ")
		bindings[i] = fmt.Sprintf("[%s [%s] %s]", e.names[i], params, e.bodies[i].str())
	}
	return fmt.Sprintf("astLetRec[%s %s]", strings.Join(bindings, " "), e.body.str())
}

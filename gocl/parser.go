package gocl

import "errors"
import "fmt"

const kw_DEF string = "def"
const kw_LET string = "let"
const kw_LETSTAR string = "let*"
const kw_LETREC string = "letrec"
const kw_LOOP string = "let"
const kw_IF string = "if"
const kw_FUN string = "fn"
const kw_QUOTE string = "quote"
const kw_DO string = "do"

const kw_MACRO string = "macro"
const kw_AND string = "and"
const kw_OR string = "or"

var fresh = (func(init int) func(string) string {
	id := init
	return func(prefix string) string {
		result := fmt.Sprintf("%s_%d", prefix, id)
		id += 1
		return result
	}
})(0)

func parseDef(sexp Value) (*astDef, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isDef := parseKeyword(kw_DEF, head)
	if !isDef {
		return nil, nil
	}
	defBlock, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to def")
	}
	if name, ok := defBlock.AsSymbol(); ok {
		// next = next.tailValue()
		head, next, ok := next.AsCons()
		if !ok {
			return nil, errors.New("too few arguments to def")
		}
		value, err := parseExpr(head)
		if err != nil {
			return nil, err
		}
		if !next.IsEmpty() {
			return nil, errors.New("too many arguments to def")
		}
		return &astDef{name, DEF_VALUE, nil, value}, nil
	}
	if head, tail, ok := defBlock.AsCons(); ok {
		name, ok := head.AsSymbol()
		if !ok {
			return nil, errors.New("definition name not a symbol")
		}
		params, err := parseSymbols(tail)
		if err != nil {
			return nil, err
		}
		//next = next.tailValue()
		head, next, ok := next.AsCons()
		if !ok {
			return nil, errors.New("too few arguments to def")
		}
		body, err := parseExpr(head)
		if err != nil {
			return nil, err
		}
		if !next.IsEmpty() {
			return nil, errors.New("too many arguments to def")
		}
		return &astDef{name, DEF_FUNCTION, params, body}, nil
	}
	return nil, errors.New("malformed def")
}

func parseExpr(sexp Value) (ast, error) {
	expr := parseAtom(sexp)
	if expr != nil {
		return expr, nil
	}
	expr, err := parseastQuote(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseastIf(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseFunction(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseLet(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseLetStar(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseastLetRec(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseDo(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	expr, err = parseastApply(sexp)
	if err != nil || expr != nil {
		return expr, err
	}
	return nil, nil
}

func parseAtom(sexp Value) ast {
	if name, ok := sexp.AsSymbol(); ok {
		return &astId{name}
	}
	if sexp.IsAtom() {
		return &astLiteral{sexp}
	}
	return nil
}

func parseKeyword(kw string, sexp Value) bool {
	name, ok := sexp.AsSymbol()
	if !ok {
		return false
	}
	return (name == kw)
}

func parseastQuote(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isQ := parseKeyword(kw_QUOTE, head)
	if !isQ {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("malformed quote")
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to quote")
	}
	return &astQuote{head1}, nil
}

func parseastIf(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isastIf := parseKeyword(kw_IF, head)
	if !isastIf {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to if")
	}
	cnd, err := parseExpr(head1)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to if")
	}
	thn, err := parseExpr(head2)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head3, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to if")
	}
	els, err := parseExpr(head3)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to if")
	}
	return &astIf{cnd, thn, els}, nil
}

func parseFunction(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, head)
	if !isFun {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to fun")
	}
	if _, ok := head1.AsSymbol(); ok {
		// we need to parse as a recursive function
		// restart from scratch
		return parseRecFunction(sexp)
	}
	params, err := parseSymbols(head1)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := parseExpr(head2)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to fun")
	}
	return makeFunction(params, body), nil
}

func parseRecFunction(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isFun := parseKeyword(kw_FUN, head)
	if !isFun {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to fun")
	}
	recName, _ := head1.AsSymbol()
	// TODO: check type of recName? guess it was already done in parseFunction...
	// next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to fun")
	}
	params, err := parseSymbols(head2)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head3, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to fun")
	}
	body, err := parseExpr(head3)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to fun")
	}
	return makeRecFunction(recName, params, body), nil
}

func parseLet(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isLet := parseKeyword(kw_LET, head)
	if !isLet {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to let")
	}
	params, bindings, err := parseBindings(head1)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to let")
	}
	body, err := parseExpr(head2)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to let")
	}
	return makeLet(params, bindings, body), nil
}

func parseLetStar(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isLet := parseKeyword(kw_LETSTAR, head)
	if !isLet {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to let*")
	}
	params, bindings, err := parseBindings(head1)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to let*")
	}
	body, err := parseExpr(head2)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to let*")
	}
	return makeLetStar(params, bindings, body), nil
}

func parseastLetRec(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isastLetRec := parseKeyword(kw_LETREC, head)
	if !isastLetRec {
		return nil, nil
	}
	//next := sexp.tailValue()
	head1, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to letrec")
	}
	names, params, bodies, err := parseFunBindings(head1)
	if err != nil {
		return nil, err
	}
	//next = next.tailValue()
	head2, next, ok := next.AsCons()
	if !ok {
		return nil, errors.New("too few arguments to letrec")
	}
	body, err := parseExpr(head2)
	if err != nil {
		return nil, err
	}
	if !next.IsEmpty() {
		return nil, errors.New("too many arguments to letrec")
	}
	return &astLetRec{names, params, bodies, body}, nil
}

func parseBindings(sexp Value) ([]string, []ast, error) {
	params := make([]string, 0)
	bindings := make([]ast, 0)
	current := sexp
	for head, next, ok := sexp.AsCons(); ok; head, next, ok = next.AsCons() {
		headB, nextB, ok := head.AsCons()
		if !ok {
			return nil, nil, errors.New("expected binding (name expr)")
		}
		name, ok := headB.AsSymbol()
		if !ok {
			return nil, nil, errors.New("expected name in binding")
		}
		params = append(params, name)
		headB2, nextB, ok := nextB.AsCons()
		if !ok {
			return nil, nil, errors.New("expected expr in binding")
		}
		if !nextB.IsEmpty() {
			return nil, nil, errors.New("too many elements in binding")
		}
		binding, err := parseExpr(headB2)
		if err != nil {
			return nil, nil, err
		}
		bindings = append(bindings, binding)
		current = next
	}
	if !current.IsEmpty() {
		return nil, nil, errors.New("malformed binding list")
	}
	return params, bindings, nil
}

func parseFunBindings(sexp Value) ([]string, [][]string, []ast, error) {
	names := make([]string, 0)
	params := make([][]string, 0)
	bodies := make([]ast, 0)
	current := sexp
	for head, next, ok := sexp.AsCons(); ok; head, next, ok = next.AsCons() {
		headB, nextB, ok := head.AsCons()
		if !ok {
			return nil, nil, nil, errors.New("expected binding (name params expr)")
		}
		name, ok := headB.AsSymbol()
		if !ok {
			return nil, nil, nil, errors.New("expected name in binding")
		}
		names = append(names, name)
		headB2, nextB, ok := nextB.AsCons()
		if !ok {
			return nil, nil, nil, errors.New("expected params in binding")
		}
		these_params, err := parseSymbols(headB2)
		if err != nil {
			return nil, nil, nil, err
		}
		params = append(params, these_params)
		headB3, nextB, ok := nextB.AsCons()
		if !ok {
			return nil, nil, nil, errors.New("expected expr in binding")
		}
		if !nextB.IsEmpty() {
			return nil, nil, nil, errors.New("too many elements in binding")
		}
		body, err := parseExpr(headB3)
		if err != nil {
			return nil, nil, nil, err
		}
		bodies = append(bodies, body)
		current = next
	}
	if !current.IsEmpty() {
		return nil, nil, nil, errors.New("malformed binding list")
	}
	return names, params, bodies, nil
}

func makeLet(params []string, bindings []ast, body ast) ast {
	return &astApply{makeFunction(params, body), bindings}
}

func makeLetStar(params []string, bindings []ast, body ast) ast {
	result := body
	for i := len(params) - 1; i >= 0; i-- {
		result = makeLet([]string{params[i]}, []ast{bindings[i]}, result)
	}
	return result
}

func makeFunction(params []string, body ast) ast {
	name := fresh("__temp")
	return &astLetRec{[]string{name}, [][]string{params}, []ast{body}, &astId{name}}
}

func makeRecFunction(recName string, params []string, body ast) ast {
	return &astLetRec{[]string{recName}, [][]string{params}, []ast{body}, &astId{recName}}
}

func parseastApply(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	fun, err := parseExpr(head)
	if err != nil {
		return nil, err
	}
	if fun == nil {
		return nil, nil
	}
	args, err := parseExprs(next)
	if err != nil {
		return nil, err
	}
	return &astApply{fun, args}, nil
}

func parseExprs(sexp Value) ([]ast, error) {
	args := make([]ast, 0)
	current := sexp
	for head, next, ok := sexp.AsCons(); ok; head, next, ok = next.AsCons() {
		curr, err := parseExpr(head)
		if err != nil {
			return nil, err
		}
		if curr == nil {
			return nil, nil
		}
		args = append(args, curr)
		current = next
	}
	if !current.IsEmpty() {
		return nil, errors.New("malformed expression list")
	}
	return args, nil
}

func parseSymbols(sexp Value) ([]string, error) {
	params := make([]string, 0)
	current := sexp
	for head, next, ok := sexp.AsCons(); ok; head, next, ok = next.AsCons() {
		name, ok := head.AsSymbol()
		if !ok {
			return nil, errors.New("expected symbol in list")
		}
		params = append(params, name)
		current = next
	}
	if !current.IsEmpty() {
		return nil, errors.New("malformed symbol list")
	}
	return params, nil
}

func parseDo(sexp Value) (ast, error) {
	head, next, ok := sexp.AsCons()
	if !ok {
		return nil, nil
	}
	isDo := parseKeyword(kw_DO, head)
	if !isDo {
		return nil, nil
	}
	exprs, err := parseExprs(next)
	if err != nil {
		return nil, err
	}
	return makeDo(exprs), nil
}

func makeDo(exprs []ast) ast {
	if len(exprs) > 0 {
		result := exprs[len(exprs)-1]
		for i := len(exprs) - 2; i >= 0; i-- {
			result = makeLet([]string{fresh("__temp")}, []ast{exprs[i]}, result)
		}
		return result
	}
	return &astLiteral{&vVoid{}}
}

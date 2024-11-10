package main

import (
	"fmt"
	"rpucella.net/go-lisp-command-language/glisp"
)

func main() {
	fmt.Println("GoLisp Command Language Standalone Interpreter 1.0.0")
	eng := glisp.NewEngine()
	eng.Repl("glisp")
}

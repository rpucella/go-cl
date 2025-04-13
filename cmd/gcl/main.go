package main

import (
	"fmt"
	"rpucella.net/go-cl/gcl"
)

func main() {
	fmt.Println("Go Command Language Standalone Interpreter 1.0.0")
	eng := gcl.NewEngine()
	eng.Repl("gcl")
}

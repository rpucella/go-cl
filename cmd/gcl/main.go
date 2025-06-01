package main

import (
	"fmt"
	"io"
	"strings"
	"bufio"
	"os"
	"rpucella.net/go-cl/gcl"
)

func main() {
	fmt.Println("Go Command Language Standalone Interpreter 1.0.0.")
	eng := gcl.NewEngine()
	eng.AddPrimitive("quit", 0, 0, primitiveQuit)
	repl(eng)
}

func repl(eng gcl.Engine) {
	prompt := "gcl"
	reader := bufio.NewReader(os.Stdin)
	direct := false
	fmt.Println("Direct mode: start line with /, toggle with isolated /.")

	for {
		tempDirect := false
		if direct {
			fmt.Printf("%s(direct)> ", prompt)
		} else {
			fmt.Printf("%s> ", prompt)
		}
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				bail()
			}
			fmt.Println("IO ERROR - ", err.Error())
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if text[0] == '/' {
			if strings.TrimSpace(text[1:]) == "" {
				direct = !direct
				continue
			} else {
				text = text[1:]
				tempDirect = true
			}
		}
		var sexp gcl.Value
		if direct || tempDirect {
			sexp, err = eng.ReadList(text)
			if err != nil {
				fmt.Println("READ ERROR -", err.Error())
				continue
			}
		} else {
			sexp, err = eng.Read(text, "")
			if err != nil {
				fmt.Println("READ ERROR -", err.Error())
				continue
			}
		}
		name, err := eng.ProcessDeclaration(sexp)
		if err != nil {
			fmt.Println("PARSE ERROR -", err.Error())
			continue
		}
		if name != "" {
			fmt.Println(name)
			continue
		}
		v, err := eng.Eval(sexp)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// need to export isNil!
		if !v.IsNil() {
			fmt.Println(v.Display())
		}
	}
}

func bail() {
	os.Exit(0)
}

func primitiveQuit(name string, args []gcl.Value) (gcl.Value, error) {
	bail()
	return gcl.NewNil(), nil
}

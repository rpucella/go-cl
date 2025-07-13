package main

import (
	"fmt"
	"io"
	"strings"
	"bufio"
	"os"
	"rpucella.net/go-cl/gocl"
)

var exitFlag bool = false

func main() {
	eng := gocl.NewEngine()
	eng.AddDefaultHelpCommand()
	eng.AddCommand("quit", "", 0, 0, nil, primitiveQuit)
	eng.AddCommand("repl", "", 0, 0, nil, mkPrimitiveRepl(eng))
	eng.AddCommand("test", "[<args> ...]", 0, -1, nil, primitiveTest)
	eng.AddPrimitive("exit", 0, 0, primitiveExit)
	if len(os.Args) > 1 {
		// We have command-line parameters.
		// Do nothing for now.
		fmt.Println("(Skip command line arguments processing)")
	}
	fmt.Println("Go Command Language Standalone Interpreter 1.0.0.")
	fmt.Println("Type help for available commands.")
	fmt.Println("Type repl for the REPL loop.")
	loop(eng)
}

func loop(eng gocl.Engine) {
	prompt := "gocli"
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s> ", prompt)
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
		sexp, err := eng.ReadCommand(text)
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %w", err))
			continue
		}
		v, err := eng.Eval(sexp)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// need to export isNil!
		if !v.IsVoid() {
			fmt.Println(v.Display())
		}
	}
}

func bail() {
	os.Exit(0)
}

func primitiveQuit(name string, args []gocl.Value) (gocl.Value, error) {
	bail()
	return gocl.NewVoid(), nil
}

func repl(eng gocl.Engine) {
	prompt := "GoLisp"
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Type (exit) to leave the repl.")
	for {
		fmt.Printf("%s> ", prompt)
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Println("IO ERROR - ", err.Error())
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		var sexp gocl.Value
		sexp, err = eng.Read(text, "")
		if err != nil {
			fmt.Println("READ ERROR -", err.Error())
			continue
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
		if exitFlag {
			// We're done!
			exitFlag = false
			return
		}
		// need to export isNil!
		if !v.IsVoid() {
			fmt.Println(v.Display())
		}
	}

}

func mkPrimitiveRepl(eng gocl.Engine) func(string, []gocl.Value) (gocl.Value, error) {
	return func(name string, args []gocl.Value) (gocl.Value, error) {
		repl(eng)
		return gocl.NewVoid(), nil
	}
}

func primitiveTest(name string, args []gocl.Value) (gocl.Value, error) {
	for _, v := range args {
		fmt.Println(v.Type(), v.Display())
	}
	return gocl.NewVoid(), nil
}

func primitiveExit(name string, args []gocl.Value) (gocl.Value, error) {
	exitFlag = true
	return gocl.NewVoid(), nil
}

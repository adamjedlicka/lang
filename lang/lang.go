package lang

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// Lang is the main structure representing the language
type Lang struct {
	hadError bool
}

// MakeLang creates new instance of the language struct
func MakeLang() Lang {
	l := Lang{}

	l.hadError = false

	return l
}

// RunFile executes source code from the file on path
func (l *Lang) RunFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	l.run(string(bytes))

	if l.hadError {
		os.Exit(65)
	}
}

// RunPrompt runs code from interactive prompt
func (l *Lang) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		l.run(text)
		l.hadError = false
	}
}

func (l *Lang) run(source string) {
	scanner := MakeScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Println(err)
		return
	}

	parser := MakeParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	interpreter := MakeInterpreter(expression)
	value, err := interpreter.Interpret()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(MakeAstPrinter().Print(expression))
	fmt.Println(interpreter.Stringify(value))
}

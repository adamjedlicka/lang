package lang

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

// Lang is the main structure representing the language
type Lang struct {
	scanner     Scanner
	parser      Parser
	interpreter Interpreter

	hadError bool
}

// MakeLang creates new instance of the language struct
func MakeLang() Lang {
	l := Lang{
		scanner:     MakeScanner(),
		parser:      MakeParser(),
		interpreter: MakeInterpreter(),
	}

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
	tokens, err := l.scanner.ScanTokens(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	stmnts, err := l.parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = l.interpreter.Interpret(stmnts)
	if err != nil {
		fmt.Println(err)
		return
	}
}

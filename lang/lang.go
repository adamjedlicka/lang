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
			panic(err)
		}

		l.run(text)
		l.hadError = false
	}
}

func (l *Lang) run(source string) {
	scanner := MakeScanner(l, source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return
	}

	parser := MakeParser(l, tokens)
	expression, err := parser.Parse()
	if err != nil {
		return
	}

	if l.hadError {
		return
	}

	fmt.Println(MakeAstPrinter().Print(expression))
}

func (l *Lang) errorSimple(line int, message string) string {
	return l.report(line, "", message)
}

func (l *Lang) errorBadToken(token Token, message string) string {
	if token.tokenType == EOF {
		return l.report(token.line, " at end", message)
	} else {
		return l.report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func (l *Lang) report(line int, where, message string) string {
	err := fmt.Sprintf("[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
	fmt.Print(err)

	return err
}

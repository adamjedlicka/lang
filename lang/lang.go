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
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.ToString())
	}
}

func (l *Lang) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lang) report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}

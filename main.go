package main

import (
	"bytes"
	"flag"
	"io/ioutil"

	"github.com/adamjedlicka/lang/src/compiler"
	"github.com/adamjedlicka/lang/src/config"
	"github.com/adamjedlicka/lang/src/vm"
)

func init() {
	flag.Parse()
}

func main() {
	if config.FlagScript != "" {
		source := loadFile(config.FlagScript)

		scanner := compiler.NewScanner(source)

		parser := compiler.NewParser(scanner)
		chunk := parser.Parse()
		if chunk == nil {
			return
		}

		vm := vm.NewVM()
		vm.Interpret(chunk)
	} else {
		flag.Usage()
	}
}

func loadFile(filename string) []rune {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return bytes.Runes(data)
}

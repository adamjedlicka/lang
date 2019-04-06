package main

import (
	"bytes"
	"flag"
	"io/ioutil"

	"github.com/adamjedlicka/lang/src/config"
	"github.com/adamjedlicka/lang/src/vm"
)

func init() {
	flag.Parse()
}

func main() {
	if config.FlagScript != "" {
		source := loadFile(config.FlagScript)

		vm := vm.NewVM()
		vm.Interpret(source)
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

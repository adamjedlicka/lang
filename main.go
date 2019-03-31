package main

import (
	"fmt"
	"os"

	"github.com/adamjedlicka/lang/lang"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lang [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		l := lang.MakeLang()
		l.RunFile(os.Args[1])
	} else {
		l := lang.MakeLang()
		l.RunPrompt()
	}
}

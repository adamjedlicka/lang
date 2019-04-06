package config

import "flag"

var (
	FlagScript string

	FlagDebug bool
	FlagStack bool
)

func init() {
	flag.StringVar(&FlagScript, "script", "", "Name of the script to be executed.")

	flag.BoolVar(&FlagDebug, "debug", false, "Enables debug messages.")
	flag.BoolVar(&FlagStack, "stack", false, "Prints out stack before every opcode.")
}

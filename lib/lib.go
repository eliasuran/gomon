package lib

import (
	"flag"
)

// lib
// functions that arent really a core part of the program
func ParseFlags() bool {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "help me")
	flag.Parse()

	return helpFlag
}

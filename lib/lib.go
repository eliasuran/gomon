package lib

import (
	"flag"
	"fmt"
	"os"
)

// lib
// functions that arent really a core part of the program

func HandleErr(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}

func ParseFlags() bool {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "help me")
	flag.Parse()

	return helpFlag
}

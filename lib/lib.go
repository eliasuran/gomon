package lib

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

// lib
// functions that arent really a core part of the program

func Starting(running <-chan bool) {
	for {
		select {
		case <-running:
			return
		default:
			color.HiBlue("Starting program...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

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

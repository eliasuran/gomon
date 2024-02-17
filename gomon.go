package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func main() {
	help := parseFlags()
	if help {
		fmt.Println("Provide the path to a directory containing a main.go file to run gomon")
		return
	}

	file := getPath()

	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print(`
 ██████╗  ██████╗ ███╗   ███╗ ██████╗ ███╗   ██╗
██╔════╝ ██╔═══██╗████╗ ████║██╔═══██╗████╗  ██║
██║  ███╗██║   ██║██╔████╔██║██║   ██║██╔██╗ ██║
██║   ██║██║   ██║██║╚██╔╝██║██║   ██║██║╚██╗██║
╚██████╔╝╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
 ╚═════╝  ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

`)
	color.HiMagenta("Made by eliasuran")
	color.Green("[ 200 ] Starting program")

	err := cmd.Run()
	handleErr("Error executing command, ", err)

}

func handleErr(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
	}
}

func parseFlags() bool {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "help me")
	flag.Parse()

	return helpFlag
}

func getPath() string {
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("No file provided")
		os.Exit(1)
	} else if len(files) != 1 {
		fmt.Println("Too many files provided")
		os.Exit(1)
	}

	return files[0]
}

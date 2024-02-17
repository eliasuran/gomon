package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

func main() {
	// sjekk help flag har blitt gitt og display hjelp message om den er det
	help := parseFlags()
	if help {
		fmt.Println("Provide the path to a directory containing a main.go file to run it using through gomon")
		return
	}

	// kul tekst
	color.HiBlue(`
 ██████╗  ██████╗ ███╗   ███╗ ██████╗ ███╗   ██╗
██╔════╝ ██╔═══██╗████╗ ████║██╔═══██╗████╗  ██║
██║  ███╗██║   ██║██╔████╔██║██║   ██║██╔██╗ ██║
██║   ██║██║   ██║██║╚██╔╝██║██║   ██║██║╚██╗██║
╚██████╔╝╚██████╔╝██║ ╚═╝ ██║╚██████╔╝██║ ╚████║
 ╚═════╝  ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ Made by eliasuran

`)

	// brukes for å displaye text som viser om programmet er i ferd med å starte eller har starta
	running := make(chan bool)
	go starting(running)

	// hent path til filen som skal kjøres
	file := getPath()

	// gjør klar til at commanden for å kjøre programmet skal kjøres
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setter sjekkingen om programmet har startet til false rett før det starter (TODO: legge til at den faktisk listener om programmet startet ordenlig)
	running <- false
	err := cmd.Run()
	handleErr("Error executing command, ", err)

	// TODO: listener som sjekker for endringer i fila
}

func starting(running <-chan bool) {
	for {
		select {
		case <-running:
			return
		default:
			color.HiBlue("Starting program.")
			time.Sleep(500 * time.Millisecond)
			color.HiBlue("Starting program..")
			time.Sleep(500 * time.Millisecond)
			color.HiBlue("Starting program...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func listener(filePath string) {
	fmt.Println(filePath)
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
		fmt.Println("No file provided\nProvide -help flag for help")
		os.Exit(1)
	} else if len(files) != 1 {
		fmt.Println("Too many files provided")
		os.Exit(1)
	}

	return files[0]
}

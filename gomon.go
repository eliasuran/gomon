package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/eliasuran/gomon/lib"
	"github.com/fatih/color"
)

func main() {
	// sjekk help flag har blitt gitt og display hjelp message om den er det
	help := lib.ParseFlags()
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
	go lib.Starting(running)

	// hent path til filen som skal kjøres
	file := getPath()

	// gjør klar til at commanden for å kjøre programmet skal kjøres
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// setter sjekkingen om programmet har startet til false rett før det starter (TODO: legge til at den faktisk listener om programmet startet ordenlig)
	running <- false

	go changeListener(file + "main.go")

	err := cmd.Run()
	lib.HandleErr("Error executing command, ", err)

	// TODO: listener som sjekker for endringer i fila
}

func changeListener(filePath string) {
	initialStat, err := os.Stat(filePath)
	lib.HandleErr("Failed to read initial file stats: ", err)

	for {
		stat, err := os.Stat(filePath)
		lib.HandleErr("Failed to reat file stats: ", err)

		if stat.Size() != initialStat.Size() {
			fmt.Println("Change in file")
		}

		time.Sleep(1 * time.Second)
	}
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

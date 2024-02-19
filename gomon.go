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
	// hent path til filen som skal kjøres
	file := getPath()

	// initialiserer command variabel
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// starter listener som sjekker for endringer i fila før programmet starter for første gang
	go changeListener(file, cmd)

	// starter programmet og sjekker for en initial error i programmet
	err := cmd.Run()
	lib.HandleErr("Error starting http server: ", err)
}

func response(status int, message string) {
	if status == 400 {
		color.HiRed("[ %d ] %s", status, message)
		return
	}

	color.HiGreen("[ %d ] %s", status, message)
}

func changeListener(filePath string, cmd *exec.Cmd) {
	response(200, "Server started successfully!")
	fullPath := filePath + "main.go"
	initialStat, err := os.Stat(fullPath)

	if err != nil {
		response(400, "Error doing something")
	} else {
		for {
			stat, err := os.Stat(fullPath)
			lib.HandleErr("Failed to read file stats: ", err)

			if stat.Size() != initialStat.Size() {
				response(200, "Change in file, resarting server...")
				initialStat = stat
			}

			time.Sleep(1 * time.Second)
		}
	}

}

func restartServer(filePath string, cmd *exec.Cmd) {
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

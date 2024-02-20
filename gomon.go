package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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

	// starter server og returnerer cmd sum bruker til å interacte ned cmd
	cmd := startServer(file)

	// initialiserer en var for prosessen til commanden som kjører
	serverProcess := cmd.Process

	// starter listener som sjekker for endringer i fila før programmet starter for første gang
	changeListener(file, serverProcess)
}

func startServer(file string) *exec.Cmd {
	// initialiserer command variabel
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// starter programmet og sjekker for en initial error i programmet
	err := cmd.Start()
	lib.HandleErr("Error starting http server: ", err)

	// gir response om at serveren har startet
	response(200, "Server started successfully!")

	// returnerer cmd
	return cmd
}

func response(status int, message string) {
	if status == 400 {
		color.HiRed("[ %d ] %s", status, message)
		return
	}

	color.HiGreen("[ %d ] %s", status, message)
}

func changeListener(filePath string, serverProcess *os.Process) {
	fullPath := filepath.Join(filePath, "main.go")
	initialStat, err := os.Stat(fullPath)

	if err != nil {
		response(400, "Error doing something")
		fmt.Println(err)
	} else {
		for {
			stat, err := os.Stat(fullPath)
			if err != nil {
				response(400, "Error in code, pls fix")
				break
			}

			if stat.ModTime() != initialStat.ModTime() {
				response(200, "Change in file, resarting server...")

				if serverProcess != nil {
					err = serverProcess.Signal(syscall.SIGTERM)
					if err != nil {
						response(400, "Error when killing previous process")
						fmt.Println(err)
					}
					serverProcess.Wait()
				}

				cmd := startServer(filePath)

				serverProcess = cmd.Process

				initialStat = stat
				response(200, "Sucessfully updated")
			}

			time.Sleep(1 * time.Second)
		}
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

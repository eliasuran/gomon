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

	// starter server og får prosessen til serveren som bruker til å interacte serveren som kjører
	process := startServer(file)

	// gir response om at serveren har startet
	response(200, "Server started successfully!")

	// starter listener som sjekker for endringer i fila før programmet starter for første gang
	changeListener(file, process)
}

func startServer(file string) *os.Process {
	// initialiserer go build cmd
	cmd := exec.Command("go", "build", "-o", filepath.Join(file, "server"), file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// kjører build commanden
	err := cmd.Run()
	if err != nil {
		response(400, "Error building program: "+err.Error())
	}

	// kjører executablen som ble laget
	cmd = exec.Command(file + "/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		response(400, "Error starting server: "+err.Error())
	}

	// lagrer prosessen i en variabel, dette gjør at jeg slipper å returnere cmd siden dette er det eneste den brukes til (for nå)
	process := cmd.Process

	// returnerer prosessen
	return process
}

func response(status int, message string) {
	t := time.Now().Format("15:04:05")
	formatting := "%s | [ %d ] %s"
	if status == 400 {
		color.Red(formatting, t, status, message)
		return
	}

	color.Green(formatting, t, status, message)
}

func changeListener(filePath string, process *os.Process) {
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

				if process != nil {
					err = process.Signal(syscall.SIGTERM)
					if err != nil {
						response(400, "Error when killing previous process")
						fmt.Println(err)
					}
					process.Wait()
				}

				process = startServer(filePath)

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

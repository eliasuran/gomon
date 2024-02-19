package main

import (
	"fmt"
	"os"
	"os/exec"
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

	// initialiserer command variabel
	cmd := exec.Command("go", "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// initialiserer process var
	serverProcess := &os.Process{Pid: -1}

	// starter listener som sjekker for endringer i fila før programmet starter for første gang
	go changeListener(file, serverProcess, cmd)

	// starter programmet og sjekker for en initial error i programmet
	err := cmd.Run()
	lib.HandleErr("Error starting http server: ", err)

	// lagrer prosessen (dette funker ikke :[)
	serverProcess = cmd.Process
}

func response(status int, message string) {
	if status == 400 {
		color.HiRed("[ %d ] %s", status, message)
		return
	}

	color.HiGreen("[ %d ] %s", status, message)
}

func changeListener(filePath string, serverProcess *os.Process, cmd *exec.Cmd) {
	response(200, "Server started successfully!")
	fullPath := filePath + "main.go"
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

				err = serverProcess.Signal(syscall.SIGINT)
				if err != nil {
					response(400, "Error when killing previous process")
					fmt.Println(err)
				}

				err = cmd.Run()
				if err != nil {
					response(400, "Error when starting up again: ")
					fmt.Println(err)
				}

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

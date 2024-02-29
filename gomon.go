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
		fmt.Println("Provide the path to a directory containing a main.go file to run it using gomon")
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
	process, err := startServer(file)
	if err != nil {
		// gir response om at en error var funnet i initial build
		response(400, "Error building server"+err.Error())
	} else {
		// gir response om at serveren har startet
		response(200, "Server started successfully!")
	}

	// starter listener som sjekker for endringer i fila før programmet starter for første gang
	changeListener(file, process)
}

func startServer(file string) (*os.Process, error) {
	// initialiserer go build cmd
	cmd := exec.Command("go", "build", "-o", filepath.Join(file, "server"), file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// kjører build commanden
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// kjører executablen som ble laget
	cmd = exec.Command(file + "/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// lagrer prosessen i en variabel, dette gjør at jeg slipper å returnere cmd siden dette er det eneste den brukes til (for nå)
	process := cmd.Process

	// returnerer prosessen
	return process, nil
}

func response(status int, message string) {
	t := time.Now().Format("15:04:05")
	formatting := "%s | [ %d ] %s"
	if status != 200 {
		color.Red(formatting, t, status, message)
		return
	}

	color.Green(formatting, t, status, message)
}

func changeListener(filePath string, process *os.Process) {
	fullPath := filepath.Join(filePath, "main.go")
	initialStat, err := os.Stat(fullPath)
	if err != nil {
		response(400, "Error getting initial stats of file: "+err.Error())
		return
	}

	for {
		stat, err := os.Stat(fullPath)
		if err != nil {
			response(400, "Error getting newest stats of file: "+err.Error())
			break
		}

		if stat.ModTime() != initialStat.ModTime() {
			response(200, "Change in file, resarting server...")

			go func() {
				process, err = startServer(filePath)
				if err != nil {
					response(400, "Error starting server: "+err.Error())
				} else {
					response(200, "Sucessfully updated")
				}
			}()

			if process != nil {
				err = process.Signal(syscall.SIGTERM)
				if err != nil {
					response(400, "Error when killing previous process: "+err.Error())
				}
				process.Wait()
			}

			initialStat = stat
		}

		time.Sleep(1 * time.Second)
	}

}

func getPath() string {
	files := os.Args[1:]
	if len(files) == 0 {
		response(400, "No file provided\nProvide -help flag for help")
		os.Exit(1)
	} else if len(files) != 1 {
		response(400, "Too many files provided\nProvide -help flag for help")
		os.Exit(1)
	}

	return files[0]
}

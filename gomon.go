package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	files := os.Args[1:]
	fmt.Println(files)
	if len(files) == 0 {
		fmt.Println("No file provided")
		os.Exit(1)
	} else if len(files) != 1 {
		fmt.Println("Too many files provided")
		os.Exit(1)
	}
	cmd := exec.Command("go", "run", files[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command. %v", err)
	}
}

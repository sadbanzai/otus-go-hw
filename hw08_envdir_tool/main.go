package main

import (
	"log"
	"os"
)

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	exitCode := RunCmd(os.Args[2:], envs)
	os.Exit(exitCode)
}

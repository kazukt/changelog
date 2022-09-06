package main

import (
	"fmt"
	"os"

	"github.com/kazukt/changelog/pkg/cmd"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := run()
	os.Exit(int(code))
}

func run() exitCode {
	rootCmd := cmd.NewCmdRoot()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		return exitError
	}

	return exitOK
}

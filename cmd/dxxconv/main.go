package main

import (
	"fmt"
	"os"

	"github.com/tetsuzawa/spat"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	in := args[0]
	out := args[1]
	data, err := spat.ReadFile(in)
	if err != nil {
		return err
	}
	return spat.WriteFile(out, data)
}

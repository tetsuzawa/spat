package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tetsuzawa/spat"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("Usage: %s signal_length(sample) out(.DXX)\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		log.Printf("error: %+v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 2 {
		return errors.New("invalid arguments")
	}
	args := flag.Args()
	samples, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	outPath := args[1]

	const fs = 48000
	pinkNoise := spat.PinkNoise(samples, fs)
	return spat.WriteDXXFile(outPath, pinkNoise)
}

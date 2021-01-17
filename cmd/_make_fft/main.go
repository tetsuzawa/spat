package main

import (
	"errors"
	"flag"
	"github.com/tetsuzawa/spat"
	"log"
	"math/cmplx"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(0)
	flag.CommandLine.Output()
	flag.Usage = func() {
		log.Printf("Usage: %s in out\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
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
	in := args[0]
	out := args[1]
	inData, err := spat.ReadDXXFile(in)
	if err != nil {
		return err
	}
	x := ToComplex(nil, inData)
	spat.FFT(x, len(x))
	if err := spat.WriteDXXFile("sine_fft_complex.DSB", Encode(nil, x)); err != nil {
		return err
	}
	spat.IFFT(x, len(x))
	//Abs(inData, x)
	//ToFloat(inData, x)
	//return spat.WriteDXXFile(out, inData)
	outData := ToFloat(nil, x)
	return spat.WriteDXXFile(out, outData)
}

func ToComplex(dst []complex128, src []float64) []complex128 {
	if dst == nil {
		dst = make([]complex128, len(src))
	}
	for i, v := range src {
		dst[i] = complex(v, 0)
	}
	return dst
}

func ToFloat(dst []float64, src []complex128) []float64 {
	if dst == nil {
		dst = make([]float64, len(src))
	}
	for i, v := range src {
		dst[i] = real(v)
	}
	return dst
}

func Abs(dst []float64, src []complex128) []float64 {
	if dst == nil {
		dst = make([]float64, len(src))
	}
	for i, v := range src {
		dst[i] = cmplx.Abs(v)
	}
	return dst
}

func Encode(dst []float64, src []complex128) []float64 {
	if dst == nil {
		dst = make([]float64, 2*len(src))
	}
	if len(dst) < len(src)*2 {
		log.Println("errrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr")
	}
	for i, v := range src {
		dst[2*i] = real(v)
		dst[2*i+1] = imag(v)
	}
	return dst
}

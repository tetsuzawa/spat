package spat

import (
	"errors"
	"log"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

// LinearConvolutionFreqDomain return linear convolution. len: len(x) + len(y) - 1
func LinearConvolutionFreqDomain(x, y []complex128) []complex128 {
	convLen := len(x) + len(y) - 1
	xPad := dsputils.ZeroPad(x, convLen)
	yPad := dsputils.ZeroPad(y, convLen)
	if len(xPad) != convLen {
		log.Fatalln("len err")
	}
	return fft.Convolve(xPad, yPad)
}

func LinearConvolutionTimeDomain(x, y []float64) []float64 {
	convLen := len(x) + len(y) - 1
	res := make([]float64, convLen)
	for p := 0; p < len(x); p++ {
		for n := p; n < len(y)+p; n++ {
			res[n] += x[p] * y[n-p]
		}
	}
	return res
}

func LinearConvolutionTimeDomainFloat32s(dst, x1, x2 []float32) ([]float32, error) {
	convLen := len(x1) + len(x2) - 1
	if dst == nil {
		dst = make([]float32, convLen)
	} else if len(dst) < convLen {
		return dst, errors.New("length of dst must be at least len(x1) + len(x2) - 1")
	}
	for i := 0; i < len(x1); i++ {
		for j := i; j < len(x2)+i; j++ {
			dst[j] += x1[i] * x2[j-i]
		}
	}
	return dst, nil
}

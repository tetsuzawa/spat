package spat

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestFFT(t *testing.T) {
	array := ToComplex(nil, makeSine(48000, 440, 48000))
	copyArray := make([]complex128, len(array))
	copy(copyArray, array)
	type args struct {
		array []T
		size  int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sine",
			args: args{
				array: array,
				size:  len(array),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FFT(tt.args.array, tt.args.size)
			IFFT(tt.args.array, tt.args.size)
			if !reflect.DeepEqual(tt.args.array, copyArray){
				fmt.Println(tt.args.array)
				t.Errorf("fft ifft return value is incorrect")
			}
		})
	}
}

func makeSine(length int, freq float64, samplingFreq int) []float64 {
	sine := make([]float64, length)
	for i := range sine {
		sine[i] = math.Sin(2 * math.Pi * freq * float64(i) / float64(samplingFreq))
	}
	return sine
}

func ToComplex(dst []complex128, x []float64) []complex128 {
	if dst == nil {
		dst = make([]complex128, len(x))
	}
	for n, v := range x {
		dst[n] = complex(v, 0)
	}
	return dst
}

func ToFloat(dst []float64, x []complex128) []float64 {
	if dst == nil {
		dst = make([]float64, len(x))
	}
	for n, v := range x {
		dst[n] = real(v)
	}
	return dst
}

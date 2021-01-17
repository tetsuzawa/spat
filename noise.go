package spat

import (
	"math"
	"math/rand"
)

// PinkNoise generates pinknoise using Voss algorithm.
// http://www.firstpr.com.au/dsp/pink-noise/
func PinkNoise(samples, fs int) []float64 {
	// lowest frequency to keep pink (Hz)
	fLow := 10
	levels := int(math.Ceil(math.Log2(float64(fs) / float64(fLow))))
	out := make([]float64, samples)
	xs := make([]float64, levels)
	for i := 0; i < samples; i++ {
		for j := 0; j < levels; j++ {
			// 1<<(j+1) is 2^(j+1)
			if i%(1<<(j+1)) == 0 {
				xs[j] = rand.NormFloat64()
			}
			out[i] = rand.NormFloat64() + Sum(xs)
		}
	}
	normFactor := MaxFloat64s(AbsFloat64s(out))
	for i, v := range out {
		out[i] = v / normFactor
	}
	return out
}

func RandNoise(samples int) []float64 {
	out := make([]float64, samples)
	for i := range out {
		out[i] = rand.Float64()
	}
	return out
}

func Sum(vs []float64) float64 {
	sum := 0.0
	for _, v := range vs {
		sum += v
	}
	return sum
}

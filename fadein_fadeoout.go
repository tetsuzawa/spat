package spat

import (
	"math"
)

func GenerateFadeinFadeoutFilter(length int) (fadeinFilter, fadeoutFilter []float64) {
	// Fourier Series Window Coefficient
	a0 := (1 + math.Sqrt(2)) / 4
	a1 := 0.25 + 0.25*math.Sqrt((5-2*math.Sqrt(2))/2)
	a2 := (1 - math.Sqrt(2)) / 4
	a3 := 0.25 - 0.25*math.Sqrt((5-2*math.Sqrt(2))/2)

	// Fourier series window
	fadeinFilter = make([]float64, length)
	fadeoutFilter = make([]float64, length)
	flength := float64(length)
	for i := 0; i < length; i++ {
		f := float64(i)
		fadeinFilter[i] = a0 - a1*math.Cos(math.Pi/flength*f) + a2*math.Cos(2.0*math.Pi/flength*f) - a3*math.Cos(3.0*math.Pi/flength*f)
		fadeoutFilter[i] = a0 + a1*math.Cos(math.Pi/flength*f) + a2*math.Cos(2.0*math.Pi/flength*f) + a3*math.Cos(3.0*math.Pi/flength*f)
	}
	return fadeinFilter, fadeoutFilter
}

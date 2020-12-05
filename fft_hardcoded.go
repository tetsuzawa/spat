package spat

const FftHardCordedSize = 4

func FFTHardCorded(array []T, size, n, sign int) {
	m := size / n
	switch m {
	case 1:
		break
	case 2:
		for i := 0; i < n; i++ {
			r0 := real(array[i*2+0]) + real(array[i*2+1])
			r1 := real(array[i*2+0]) + (-real(array[i*2+1]))
			i0 := imag(array[i*2+0]) + imag(array[i*2+1])
			i1 := imag(array[i*2+0]) + (-imag(array[i*2+1]))
			array[i*2+0] = complex(r0, i0)
			array[i*2+1] = complex(r1, i1)
		}
		break
	case 3:
		for i := 0; i < n; i++ {
			r0 := real(array[i*3+0]) + real(array[i*3+1]) + real(array[i*3+2])
			r1 := real(array[i*3+0]) + (-5.00000000000000e-01)*(real(array[i*3+1])) + (-5.00000000000000e-01)*(real(array[i*3+2])) + (-float64(sign))*(8.66025403784439e-01)*(imag(array[i*3+1])) + (-float64(sign))*(-8.66025403784438e-01)*(imag(array[i*3+2]))
			r2 := real(array[i*3+0]) + (-5.00000000000000e-01)*(real(array[i*3+1])) + (-5.00000000000000e-01)*(real(array[i*3+2])) + (-float64(sign))*(-8.66025403784438e-01)*(imag(array[i*3+1])) + (-float64(sign))*(8.66025403784439e-01)*(imag(array[i*3+2]))
			i0 := imag(array[i*3+0]) + (imag(array[i*3+1])) + (imag(array[i*3+2]))
			i1 := float64(sign)*(8.66025403784439e-01)*(real(array[i*3+1])) + float64(sign)*(-8.66025403784438e-01)*(real(array[i*3+2])) + (imag(array[i*3+0])) + (-5.00000000000000e-01)*(imag(array[i*3+1])) + (-5.00000000000000e-01)*(imag(array[i*3+2]))
			i2 := float64(sign)*(-8.66025403784438e-01)*(real(array[i*3+1])) + float64(sign)*(8.66025403784439e-01)*(real(array[i*3+2])) + (imag(array[i*3+0])) + (-5.00000000000000e-01)*(imag(array[i*3+1])) + (-5.00000000000000e-01)*(imag(array[i*3+2]))
			array[i*3+0] = complex(r0, i0)
			array[i*3+1] = complex(r1, i1)
			array[i*3+2] = complex(r2, i2)
		}
		break
	case 4:
		for i := 0; i < n; i++ {
			r0 := real(array[i*4+0]) + real(array[i*4+1]) + real(array[i*4+2]) + real(array[i*4+3])
			r1 := real(array[i*4+0]) + (-real(array[i*4+2])) + (-float64(sign))*(imag(array[i*4+1])) + float64(sign)*(imag(array[i*4+3]))
			r2 := real(array[i*4+0]) + (-real(array[i*4+1])) + real(array[i*4+2]) + (-real(array[i*4+3]))
			r3 := real(array[i*4+0]) + (-real(array[i*4+2])) + float64(sign)*(imag(array[i*4+1])) + (-float64(sign))*(imag(array[i*4+3]))
			i0 := imag(array[i*4+0]) + (imag(array[i*4+1])) + (imag(array[i*4+2])) + (imag(array[i*4+3]))
			i1 := float64(sign)*(real(array[i*4+1])) + (-float64(sign))*(real(array[i*4+3])) + (imag(array[i*4+0])) + (-imag(array[i*4+2]))
			i2 := imag(array[i*4+0]) + (-imag(array[i*4+1])) + (imag(array[i*4+2])) + (-imag(array[i*4+3]))
			i3 := (-float64(sign))*(real(array[i*4+1])) + float64(sign)*(real(array[i*4+3])) + (imag(array[i*4+0])) + (-imag(array[i*4+2]))
			array[i*4+0] = complex(r0, i0)
			array[i*4+1] = complex(r1, i1)
			array[i*4+2] = complex(r2, i2)
			array[i*4+3] = complex(r3, i3)
		}
		break
	}
}

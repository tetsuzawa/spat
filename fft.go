package spat

import "math"

const FFTPrimeThreshold = 256

type T = complex128

func SizedDFT(array, temp, omega []complex128, size, n int) {
	m := size / n
	for i := 0; i < size; i++ {
		temp[i] = 0
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			for k := 0; k < m; k++ {
				temp[i*m+j] += array[i*m+k] * omega[((j*k)%m)*n]
			}
		}
	}
	for i := 0; i < size; i++ {
		array[i] = temp[i]
	}
}

func DFT(array []complex128, size int) {
	temp := make([]complex128, size)
	omega := make([]complex128, size)
	for i := 0; i < size; i++ {
		omega[i] = complex(FFTCos(i, len(omega)), -FFTSin(i, len(omega)))
	}
	SizedDFT(array, temp, omega, size, 1)
}

func IDFT(array []complex128, size int) {
	temp := make([]complex128, size)
	omega := make([]complex128, size)
	for i := 0; i < size; i++ {
		omega[i] = complex(FFTCos(i, len(omega)), FFTSin(i, len(omega)))
	}
	SizedDFT(array, temp, omega, size, -1)
	for i := 0; i < size; i++ {
		array[i] = complex(real(array[i])/float64(size), imag(array[i])/float64(size))
	}
}

func SizedFFT2(array []T, temp, omega []complex128, size int) {
	n := 1
	for ; n < size; n <<= 1 {
		m := size / n
		q := m >> 1
		/*permutation*/
		for i := 0; i < n; i++ {
			for k := 0; k < q; k++ {
				temp[i*m+k] = array[i*m+2*k]
				temp[i*m+q+k] = array[i*m+2*k+1]
			}
		}
		for i := 0; i < size; i++ {
			array[i] = temp[i]
		}
	}
	for ; n > 1; n >>= 1 {
		m := size / n
		q := n >> 1
		r := size / q
		/*adding up with twiddle factors*/
		for i := 0; i < size; i++ {
			temp[i] = 0
		}
		for i := 0; i < q; i++ {
			for k := 0; k < m; k++ {
				temp[2*m*i+k] += array[2*m*i+k]
				temp[2*m*i+m+k] += array[2*m*i+k]
				temp[2*m*i+k] += array[2*m*i+m+k] * omega[(k%r)*q]
				temp[2*m*i+m+k] += array[2*m*i+m+k] * omega[((m+k)%r)*q]
			}
		}
		for i := 0; i < size; i++ {
			array[i] = temp[i]
		}
	}
}

func FFT2(array, temp, omega []complex128) {
	SizedFFT2(array, temp, omega, len(array))
}

func FFTPrime(array []T, size, n, sign int) {
	p := size / n
	length := p - 1
	power2 := true
	i := 0
	for ; length > 1; length >>= 1 {
		if length&1 != 0 {
			power2 = false
		}
		i++
	}
	if power2 {
		length = p - 1
	} else {
		length = 1 << (i + 2)
	}
	pad := length - (p - 1)
	g := primePrimitiveRoot(p)
	ig := powerMod(g, p-2, p)
	data := make([]complex128, length)
	temp := make([]complex128, length)
	omega := make([]complex128, length)
	omega2 := make([]complex128, length)

	for i := 0; i < length; i++ {
		omega[i] = complex(FFTCos(i, length), -FFTSin(i, length))
		omega2[i] = complex(FFTCos(powerMod(ig, i%(p-1), p), p), float64(sign)*FFTSin(powerMod(ig, i%(p-1), p), p))
	}
	FFT2(omega2, temp, omega)
	for i := 0; i < n; i++ {
		data[0] = array[i*p+1]
		for j := 1; j < length; j++ {
			if j <= pad {
				data[j] = 0
			} else {
				data[j] = array[i*p+powerMod(g, j-pad, p)]
			}
		}
		/*===Convolution theorem===*/
		FFT2(data, temp, omega)
		for j := 0; j < length; j++ {
			temp[j] = data[j] * omega2[j]
		}
		for j := 0; j < length; j++ {
			data[j] = temp[j]
			omega[i] = complex(real(omega[i]), imag(-omega[i]))

		}
		FFT2(data, temp, omega)
		/*add array[i*p] term*/
		for j := 0; j < length; j++ {
			data[j] = complex(real(data[j])/float64(length)+real(array[i*p]), imag(data[j])/float64(length)+imag(array[i*p]))

		}
		/*===Convolution theorem end===*/

		/*DC term*/
		temp[0] = 0
		for j := 0; j < p; j++ {
			temp[0] += array[i*p+j]
		}
		array[i*p] = temp[0]

		/*non-DC term*/
		for j := 0; j < p-1; j++ {
			array[i*p+powerMod(ig, j, p)] = data[j]
		}
	}
}

func FFT(array []T, size, sign int) {
	temp := make([]complex128, size)
	omega := make([]complex128, size)
	for i := 0; i < size; i++ {
		omega[i] = complex(FFTCos(i, size), float64(sign)*FFTSin(i, size))
	}
	var primes [32]int
	np := 0
	n := 1
	for ; ; {
		m := size / n
		p := 2 + (m & 1)
		for ; p*p <= m; p += 2 {
			if m%p == 0 {
				/*m is divisible by p*/
				q := m / p
				/*permutation*/
				for i := 0; i < n; i++ {
					for j := 0; j < p; j++ {
						for k := 0; k < q; k++ {
							temp[i*m+q*j+k] = array[i*m+p*k+j]
						}
					}
				}
				for i := 0; i < size; i++ {
					array[i] = temp[i]
				}
				primes[np] = p
				np++
				n *= p
				if q <= FftHardCordedSize {
					break
				}
				goto next
			}
		}
		break
	next:
	}
	/*bottom of recursion*/
	/*perform dft on n sub-arrays*/
	{
		if size/n <= FftHardCordedSize {
			FFTHardCorded(array, size, n, sign)
		} else if size/n <= FFTPrimeThreshold {
			SizedDFT(array, temp, omega, size, n)
		} else {
			FFTPrime(array, size, n, sign)
		}
	}
	/*sum up with twiddle factors*/
	for h := np - 1; h >= 0; h-- {
		m := size / n
		p := primes[h]
		q := n / p
		r := size / q

		for i := 0; i < size; i++ {
			temp[i] = 0
		}
		for i := 0; i < q; i++ {
			for j := 0; j < p; j++ {
				for k := 0; k < m; k++ {
					for l := 0; l < p; l++ {
						temp[i*p*m+(j*m+k)] +=
							array[i*p*m+(l*m+k)] * omega[(l*(j*m+k)%r)*q]
					}
				}
			}
		}
		for i := 0; i < size; i++ {
			array[i] = temp[i]
		}
		n = q
	}
}

func FFTCos(n, m int) float64 {
	if n == 0 {
		return 1
	}
	if n == m { /*2 Pi*/ return 1
	}
	if m%n == 0 {
		if m/n == 2 { /*Pi*/ return -1
		}
		if m/n == 4 { /*Pi/2*/ return 0
		}
	}
	if m%(m-n) == 0 {
		if m/(m-n) == 4 { /*3 Pi/2*/ return 0
		}
	}
	return math.Cos(2 * math.Pi * float64(n) / float64(m))
}

func FFTSin(n, m int) float64 {
	if n == 0 {
		return 0
	}
	if n == m { /*2 Pi*/ return 0
	}
	if m%n == 0 {
		if m/n == 2 { /*Pi*/ return 0
		}
		if m/n == 4 { /*Pi/2*/ return 1
		}
	}
	if m%(m-n) == 0 {
		if m/(m-n) == 4 { /*3 Pi/2*/ return -1
		}
	}
	return math.Sin(2 * math.Pi * float64(n) / float64(m))
}

func primePrimitiveRoot(p int) int {
	var i, t int
	var list []int

	t = p - 1
	for i = 2; i*i <= t; i++ {
		if t%i == 0 {
			list = append(list, i)
			t /= i
			for t%i == 0 {
				t /= i
			}
		}
	}

	for idx, v := range list {
		list[idx] = (p - 1) / v
	}
	for i = 2; i <= p-1; i++ {
		for _, v := range list {
			if powerMod(i, v, p) == 1 {
				goto loopend
			}
		}
		break
	loopend:
	}
	return i
}

func powerMod(a, b, m int) int {
	var i int
	for i = 1; b != 0; b >>= 1 {
		if b&1 != 0 {
			i = (i * a) % m
		}
		a = (a * a) % m
	}
	return i
}

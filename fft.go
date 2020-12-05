package spat

type T = complex128

func SizedDFT(array, temp, omega []complex128, size, n int) {
	var m int = size / n
	for i := 0; i < size; i++ {
		temp[i] = 0;
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			for k := 0; k < m; k++ {
				temp[i*m+j] += array[i*m+k] * omega[((j*k)%m)*n];
			}
		}
	}
	for i := 0; i < size; i++ {
		array[i] = temp[i];
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
		array[i] = complex(real(array[i])/float64(size), imag(array[i])/float64(size), );
	}
}

func fft2(array []T, temp, omega []complex128, size int) {
	n := 1
	for ; n < size; n <<= 1 {
		m := size / n;
		q := m >> 1;
		/*permutation*/
		for i := 0; i < n; i++ {
			for k := 0; k < q; k++ {
				temp[i*m+k] = array[i*m+2*k]
				temp[i*m+q+k] = array[i*m+2*k+1]
			}
		}
		for i := 0; i < size; i++ {
			array[i] = temp[i];
		}
	}
	for ; n > 1; n >>= 1 {
		m := size / n;
		q := n >> 1;
		r := size / q;
		/*adding up with twiddle factors*/
		for i := 0; i < size; i++ {
			temp[i] = 0;
		}
		for i := 0; i < q; i++ {
			for k := 0; k < m; k++ {
				temp[2*m*i+k] += array[2*m*i+k];
				temp[2*m*i+m+k] += array[2*m*i+k];
				temp[2*m*i+k] += array[2*m*i+m+k] * omega[(k%r)*q];
				temp[2*m*i+m+k] += array[2*m*i+m+k] * omega[((m+k)%r)*q];
			}
		}
		for i := 0; i < size; i++ {
			array[i] = temp[i];
		}
	}
}

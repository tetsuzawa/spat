package spat

type T float64

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

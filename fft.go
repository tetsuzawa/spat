package spat

func DFT(array, temp, omega []complex128, size, n int) {
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

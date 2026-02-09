package matrix

// MatrixMultiply multiplies two matrices a and b.
func Multiply(a, b [][]float64) [][]float64 {
	out := make([][]float64, len(a))
	for i := range a {
		out[i] = make([]float64, len(b[0]))
		for j := range b[0] {
			for k := range b {
				out[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return out
}

// MatrixTranspose returns the transpose of a matrix m.
func Transpose(m [][]float64) [][]float64 {
	r, c := len(m), len(m[0])
	out := make([][]float64, c)
	for i := range out {
		out[i] = make([]float64, r)
		for j := range m {
			out[i][j] = m[j][i]
		}
	}
	return out
}

// GaussianElimination solves the system of linear equations Ax = b using Gaussian elimination.
func GaussianElimination(a [][]float64, b []float64) []float64 {
	n := len(b)

	for i := 0; i < n; i++ {
		pivot := a[i][i]
		for j := i; j < n; j++ {
			a[i][j] /= pivot
		}
		b[i] /= pivot

		for k := 0; k < n; k++ {
			if k == i {
				continue
			}
			factor := a[k][i]
			for j := i; j < n; j++ {
				a[k][j] -= factor * a[i][j]
			}
			b[k] -= factor * b[i]
		}
	}

	return b
}

// SolveLeastSquares solves the least squares problem Ax = b.
func SolveLeastSquares(a [][]float64, b []float64) []float64 {
	AT := Transpose(a)
	ATA := Multiply(AT, a)

	ATb := make([]float64, len(AT))
	for i := range AT {
		for j := range b {
			ATb[i] += AT[i][j] * b[j]
		}
	}

	return GaussianElimination(ATA, ATb)
}

// AddRegularization adds λI to the matrix a.
func AddRegularization(a [][]float64, lambda float64) {
	for i := range a {
		a[i][i] += lambda
	}
}

// SolveLeastSquaresRegularized solves the regularized least squares problem (A^T A + λI)x = A^T b.
func SolveLeastSquaresRegularized(a [][]float64, b []float64, lambda float64) []float64 {
	AT := Transpose(a)
	ATA := Multiply(AT, a)

	// Add λI
	AddRegularization(ATA, lambda)

	ATb := make([]float64, len(AT))
	for i := range AT {
		for j := range b {
			ATb[i] += AT[i][j] * b[j]
		}
	}

	return GaussianElimination(ATA, ATb)
}

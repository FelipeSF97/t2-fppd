package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	N    = 3000
	Seed = 42
)

func createMatrix(n int) []float64 {
	return make([]float64, n*n)
}

func fillMatrix(m []float64) {
	for i := range m {
		m[i] = rand.Float64() * 10
	}
}

func multiplyMatrices(a, b, c []float64, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += a[i*n+k] * b[k*n+j]
			}
			c[i*n+j] = sum
		}
	}
}

func printVerification(c []float64, n int) {
	fmt.Println("\nValores para verificação:")
	fmt.Printf("C[0][0] = %.2f\n", c[0])
	fmt.Printf("C[0][N-1] = %.2f\n", c[n-1])
	fmt.Printf("C[N-1][0] = %.2f\n", c[(n-1)*n])
	fmt.Printf("C[N-1][N-1] = %.2f\n", c[n*n-1])

	checksum := 0.0
	for _, v := range c {
		checksum += v
	}
	fmt.Printf("Checksum = %.2f\n", checksum)
}

func main() {
	rand.Seed(Seed)

	fmt.Printf("Gerando matrizes %dx%d...\n", N, N)

	A := createMatrix(N)
	B := createMatrix(N)
	C := createMatrix(N)

	fillMatrix(A)
	fillMatrix(B)

	fmt.Println("Iniciando multiplicação...")

	start := time.Now()

	multiplyMatrices(A, B, C, N)

	elapsed := time.Since(start)

	fmt.Printf("\nTempo de execução: %v\n", elapsed)

	printVerification(C, N)
}

package main

import (
	"fmt"
	"math/rand"

	mpi "github.com/mvneves/gompi"
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

func getWorkload(rank, size, n int) (startRow, endRow int) {
	rowsPerProcess := n / size
	remainder := n % size

	startRow = rank * rowsPerProcess
	endRow = startRow + rowsPerProcess

	if rank == size-1 {
		endRow += remainder
	}

	return
}

func multiplyLocal(localA, localB, localC []float64, numRows, n int) {
	for i := 0; i < numRows; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0

			for k := 0; k < n; k++ {
				sum += localA[i*n+k] * localB[k*n+j]
			}

			localC[i*n+j] = sum
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
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)

	rank := world.GetRank()
	size := world.GetSize()

	fmt.Printf("Processo %d de %d iniciado.\n", rank, size)

	startRow, endRow := getWorkload(rank, size, N)

	numRows := endRow - startRow

	localA := make([]float64, numRows*N)
	localB := createMatrix(N)
	localC := make([]float64, numRows*N)

	fmt.Printf(
		"Processo %d receberá %d linhas (%d elementos).\n",
		rank,
		numRows,
		len(localA),
	)

	fmt.Printf(
		"Processo %d ficará responsável pelas linhas %d até %d\n",
		rank,
		startRow,
		endRow-1,
	)

	if rank == 0 {
		rand.Seed(Seed)

		fmt.Println("Gerando matrizes...")

		A := createMatrix(N)
		B := createMatrix(N)
		C := createMatrix(N)

		fillMatrix(A)
		fillMatrix(B)

		for dest := 1; dest < size; dest++ {

			start, end := getWorkload(dest, size, N)

			world.Send(
				A[start*N:end*N],
				dest,
				0,
			)
		}

		copy(
			localA,
			A[startRow*N:endRow*N],
		)

		copy(localB, B)

		for dest := 1; dest < size; dest++ {
			world.Send(B, dest, 1)
		}

		fmt.Printf("Matriz A: %d elementos\n", len(A))
		fmt.Printf("Matriz B: %d elementos\n", len(B))
		fmt.Printf("Matriz C: %d elementos\n", len(C))

		multiplyLocal(localA, localB, localC, numRows, N)

		fmt.Printf(
			"Processo %d terminou seu bloco.\n",
			rank,
		)

		copy(
			C[startRow*N:endRow*N],
			localC,
		)

		for src := 1; src < size; src++ {

			start, end := getWorkload(src, size, N)

			buffer := make([]float64, (end-start)*N)

			world.Recv(
				&buffer,
				src,
				2,
			)

			copy(
				C[start*N:end*N],
				buffer,
			)
		}

		fmt.Println("Todos os blocos foram recebidos.")

		printVerification(C, N)
	} else {

		world.Recv(
			&localA,
			0,
			0,
		)

		world.Recv(&localB, 0, 1)

		fmt.Printf(
			"Processo %d recebeu %d elementos de B\n",
			rank,
			len(localB),
		)

		fmt.Printf(
			"Processo %d recebeu %d elementos de A\n",
			rank,
			len(localA),
		)

		multiplyLocal(localA, localB, localC, numRows, N)

		fmt.Printf(
			"Processo %d terminou seu bloco.\n",
			rank,
		)

		world.Send(
			localC,
			0,
			2,
		)
	}
}

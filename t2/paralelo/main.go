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

	// O último processo recebe as linhas restantes
	if rank == size-1 {
		endRow += remainder
	}

	return
}

func main() {
	mpi.Init()
	defer mpi.Finalize()

	world := mpi.NewComm(true)

	rank := world.GetRank()
	size := world.GetSize()

	startRow, endRow := getWorkload(rank, size, N)

	numRows := endRow - startRow

	localA := createMatrix(numRows)
	//localC := createMatrix(numRows)

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
		// Seed fixa
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

		fmt.Printf("Matriz A: %d elementos\n", len(A))
		fmt.Printf("Matriz B: %d elementos\n", len(B))
		fmt.Printf("Matriz C: %d elementos\n", len(C))
	} else {

		world.Recv(
			&localA,
			0,
			0,
		)

		fmt.Printf(
			"Processo %d recebeu %d elementos de A\n",
			rank,
			len(localA),
		)

	}

	fmt.Printf("Processo %d de %d iniciado.\n", rank, size)
}

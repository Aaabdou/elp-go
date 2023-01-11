package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	matrixA, err := readMatrixFromFile("matrix_a.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	matrixB, err := readMatrixFromFile("matrix_b.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := matrixMultiply(matrixA, matrixB)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, row := range result {
		for _, value := range row {
			fmt.Printf("%d ", value)
		}
		fmt.Println()
	}
}


// The function returns a 2D slice representing the matrix, or an error if the file could not be read.
func readMatrixFromFile(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matrix [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowStrings := strings.Fields(scanner.Text())
		var row []int
		for _, valueString := range rowStrings {      // range iterates over the slice elements 
			value, err := strconv.Atoi(valueString)   // convert string to int 
			if err != nil {
				return nil, err
			}
			row = append(row, value)                  // apprend new elements to the slice
		}
		matrix = append(matrix, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	//fmt.Println(matrix)
	
	return matrix, nil
}


func matrixMultiply(matrixA, matrixB [][]int) ([][]int, error) {
	if len(matrixA[0]) != len(matrixB) {
		return nil, fmt.Errorf("Number of columns in matrix A does not match number of rows in matrix B")
	}

	result := make([][]int, len(matrixA))    // create a 2D slice with the size = the length of matrix A
	for i := range result {                  // initialize a slice with length of matrix B in each of slices
		result[i] = make([]int, len(matrixB[0]))
	}

	for i := 0; i < len(matrixA); i++ {
		for j := 0; j < len(matrixB[0]); j++ {
			for k := 0; k < len(matrixB); k++ {
				result[i][j] += matrixA[i][k] * matrixB[k][j]
			}
		}
	}
	//fmt.Println(result)
	return result, nil
}

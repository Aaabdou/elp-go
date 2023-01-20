package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	//"net"
	"sync"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type line struct {
	nb   int
	line []int
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	var matrixA [][]int
	var matrixB [][]int
	wg.Add(2)
	go func() {
		var errA error
		matrixA, errA = readMatrixFromFile("matrix_a.txt")
		if errA != nil {
			fmt.Errorf("Can't read Matrix A")
		}
		wg.Done()
	}()
	go func() {

		var errB error
		matrixB, errB = readMatrixFromFile("matrix_b.txt")
		if errB != nil {
			fmt.Errorf("Can't read Matrix B")
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println(matrixA)
	fmt.Println(matrixB)

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
	fmt.Println("Exe time =", time.Since(start))
}

// The function returns a 2D slice representing the matrix, or an error if the file could not be read.
func readMatrixFromFile(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() //close file at func end

	//mLength, err := lineCounter(file)
	scanner := bufio.NewScanner(file)
	var mLength = 0
	can := make(chan line)
	for scanner.Scan() {
		rowStrings := strings.Fields(scanner.Text()) // Sépare les strings en slice sans espace
		go readLine(can, rowStrings, mLength)
		mLength++
	}
	var matrix [][]int = make([][]int, mLength)
	for j := 0; j < mLength; j++ {
		mLine := <-can
		matrix[mLine.nb] = mLine.line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	//fmt.Println(matrix)

	return matrix, nil

}
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
func readLine(c chan line, rowStrings []string, i int) error {
	var row []int
	l := line{nb: i}
	for _, valueString := range rowStrings { // range iterates over the slice elements
		value, err := strconv.Atoi(valueString) // convert string to int
		if err != nil {
			return err
		}
		row = append(row, value) // apprend new elements to the slice
	}
	l.line = row
	c <- l
	//fmt.Println(l.nb, " : ", l.line)
	return nil
}

func matrixMultiply(matrixA, matrixB [][]int) ([][]int, error) {
	if len(matrixA[0]) != len(matrixB) {
		return nil, fmt.Errorf("Number of columns in matrix A does not match number of rows in matrix B")
	}

	result := make([][]int, len(matrixA)) // create a 2D slice with the size = the length of matrix A

	can := make(chan line)
	for i := 0; i < len(matrixA); i++ {
		go lineXMatrix(matrixA[i], matrixB, can, i)
	}
	for j := 0; j < len(matrixA); j++ {
		mLine := <-can
		result[mLine.nb] = mLine.line
	}

	//fmt.Println(result)
	return result, nil
}
func lineXMatrix(lin []int, matrix [][]int, c chan line, index int) {
	var result = make([]int, len(matrix[0]))
	l := line{nb: index}
	for j := 0; j < len(matrix[0]); j++ {
		for k := 0; k < len(lin); k++ {
			result[j] += lin[k] * matrix[k][j]
		}
	}
	l.line = result
	c <- l
}

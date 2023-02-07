package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	SERV_HOST = "localhost"
	SERV_PORT = "3333"
	SERV_TYPE = "tcp"
)

type line struct {
	nb   int
	line []int
	err  error
}

func main() {
	wg := new(sync.WaitGroup)
	chConn := make(chan net.Conn)
	running := true
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(chConn, wg)
	}
	fmt.Println("Server starting...")
	l, err := net.Listen(SERV_TYPE, SERV_HOST+":"+SERV_PORT)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Listen for an incoming connection.
	go user(&running)
	for running {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		chConn <- conn

	}
}
func user(run *bool) {
	fmt.Println("If you want to close the server, type 1")
	var answer string
	for true {
		fmt.Scanln(&answer)
		if answer == "1" {
			*run = false
			break
		} else {
			fmt.Println("Didn't understand your request")
			fmt.Println("If you want to close the server, type 1")
		}
	}

}
func worker(can chan net.Conn, wg *sync.WaitGroup) {
	client_exit := false
	for conn := range can {
		fmt.Println("Nouveau client")
		for !client_exit {

			var Aexist = false
			var Bexist = false
			var matrixA [][]int
			var matrixB [][]int

			for !Aexist || !Bexist {
				rd := bufio.NewReader(conn)
				ba, err := rd.ReadString(byte(':'))
				name := strings.Trim(ba, ":")
				fmt.Println(name)
				if err != nil {
					fmt.Println("Error encountered", err)
					client_exit = true
					break
				}

				if name == "matrixA" {
					ba, err = rd.ReadString(byte(':'))
					if err != nil {
						fmt.Println("Error encountered", err)
						client_exit = true
						break
					}
					ba = strings.Trim(ba, ":")
					fileSize, _ := strconv.ParseInt(ba, 10, 64)
					fmt.Println(name, " size :", fileSize)
					matrixA, err = readMatrixFromClient(rd)
					if err != nil {

						fmt.Println("Error encountered", err)
						continue
					}
					Aexist = true

				} else if name == "matrixB" {
					ba, err = rd.ReadString(byte(':'))
					if err != nil {
						fmt.Println("Error encountered", err)
						client_exit = true
						break
					}
					ba = strings.Trim(ba, ":")
					fileSize, _ := strconv.ParseInt(ba, 10, 64)
					fmt.Println(name, " size :", fileSize)
					matrixB, err = readMatrixFromClient(rd)
					if err != nil {
						fmt.Println("Error encountered", err)
						continue
					}
					Bexist = true

				}
			}
			if !client_exit {
				fmt.Println(matrixA)
				fmt.Println(matrixB)
				result, err := matrixMultiply(matrixA, matrixB)
				if err != nil {
					fmt.Println("Error encountered", err)
					return
				}
				sendMatrix("matrixC", result, conn)
			}
		}
	}
}

func readMatrixFromClient(file io.Reader) ([][]int, error) {
	scanner := bufio.NewScanner(file)
	var mLength = 0
	can := make(chan line)

	for scanner.Scan() {
		//fmt.Println("test", mLength)
		if scanner.Text() == "EOF" {
			break
		}
		rowStrings := strings.Fields(scanner.Text()) // Sépare les strings en slice sans espace

		go readLine(can, rowStrings, mLength)
		mLength++

	}
	b := scanner.Err()
	if b != nil {
		fmt.Println(b)
	}

	var matrix [][]int = make([][]int, mLength)
	for j := 0; j < mLength; j++ {
		//fmt.Println("test", mLength)
		mLine := <-can
		if mLine.err != nil {
			return matrix, mLine.err
		}
		matrix[mLine.nb] = mLine.line
	}

	//fmt.Println(matrix)

	return matrix, nil
}
func readLine(c chan line, rowStrings []string, i int) {
	var row []int
	l := line{nb: i}
	l.err = nil
	for _, valueString := range rowStrings { // range iterates over the slice elements
		//fmt.Println(i, ":", t, "_", valueString)
		value, err := strconv.Atoi(valueString) // convert string to int
		//fmt.Println(i, ":", t, "_", value)
		if err != nil {
			l.err = err
			fmt.Println(err)
		}
		row = append(row, value) // apprend new elements to the slice
	}
	//fmt.Println("row ", i, " ", row)
	l.line = row
	fmt.Println(l.nb, " : ", l.line, ":", l.err)
	c <- l

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
func sendMatrix(name string, m [][]int, conn net.Conn) error {

	//Envoi le nom
	fmt.Println(name)

	wr := bufio.NewWriter(conn)
	wr.WriteString(name + ":")
	fmt.Println(len(m), ":", len(m[0]))
	wr.WriteRune(int32(len(m)))
	//wr.WriteString(":")
	wr.WriteRune(int32(len(m[0])))

	wr.Flush()

	for _, row := range m {
		for _, value := range row {
			fmt.Printf("%d ", int32(value))
			_, err := wr.WriteRune(int32(value))
			if err != nil {
				fmt.Println("error:", err)
			}
		}
		fmt.Println()
	}
	wr.WriteString("EOF\n")

	wr.Flush()

	return nil

}

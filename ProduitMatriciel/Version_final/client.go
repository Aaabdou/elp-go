package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	conn, err := net.Dial("tcp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Couldn't connect to server. Error :")
		fmt.Println(err)
		return
	}
	fmt.Println("Tcp connection established on port", CONN_PORT)
	for {
		fmt.Println("1. Send problem from file to server")
		fmt.Println("2. Send problem from generated matrix to server")
		fmt.Println("3. End connection")
		fmt.Println()
		var answer string
		fmt.Scanln(&answer)
		if answer == "1" {
			s, err := client(conn, true)
			if s == "Error connection" {
				fmt.Println("Error encountered :", err)
				break
			} else if err != nil {
				fmt.Println("Error encountered :", err)
			}
			if s == "STOP" {
				fmt.Println("TCP client exiting...")
				/*wr := bufio.NewWriter(conn)
				wr.WriteString("STOP:")*/
				conn.Close()
				break
			}

		} else if answer == "2" {
			s, err := client(conn, false)
			if s == "Error connection" {
				fmt.Println("Error encountered :", err)
				break
			} else if err != nil {
				fmt.Println("Error encountered :", err)
			}
			if s == "STOP" {
				fmt.Println("TCP client exiting...")
				wr := bufio.NewWriter(conn)
				wr.WriteString("STOP:")
				conn.Close()
				break
			}

		} else if answer == "3" {
			conn.Close()
			break
		} else {
			fmt.Println("Couldn't understand your answer. Please choose between")
		}

	}

}
func client(conn net.Conn, file bool) (string, error) {
	fmt.Println("You can end connection anytime by typing STOP")
	var result [][]int
	var matA string
	var colA string
	var matB string

	//Demander le fichier des matrices
	if file {
		fmt.Print("Adress matrixA:")
		fmt.Scanln(&matA)
		if matA == "STOP" {
			return "STOP", nil
		}

		fmt.Print("Adress matrixB:")
		fmt.Scanln(&matB)
		if matB == "STOP" {
			return "STOP", nil
		}

		//Envoi des matrices au serveur
		err := sendFile("matrixA", matA, conn)
		if err != nil {
			return "Error", err
		}
		err = sendFile("matrixB", matB, conn)
		if err != nil {
			return "Error", err
		}
	} else {
		var rowA int
		var columnA int
		var columnB int
		var err error
		for true {

			fmt.Println("Rows matrixA:")
			fmt.Scanln(&matA)
			if matA == "STOP" {
				return "STOP", nil
			}
			if matA == "random" {
				rowA = 0
				break
			} else {
				rowA, err = strconv.Atoi(matA)
				if err == nil {
					break
				} else {
					fmt.Println("Couldn't understand your input :(")
				}
			}
		}
		for true {
			fmt.Println("Columns matrixA:")
			fmt.Scanln(&colA)
			if colA == "STOP" {
				return "STOP", nil
			}
			if colA == "random" {
				columnA = 0
				break
			} else {
				columnA, err = strconv.Atoi(colA)
				if err == nil {
					break
				} else {
					fmt.Println("Couldn't understand your input :(")
				}
			}
		}

		matrixA, _, columnA := matrixGenerator(rowA, columnA)
		fmt.Println(matrixA)
		for true {
			var err error
			fmt.Print("Columns matrixB:")
			fmt.Scanln(&matB)

			if matB == "STOP" {
				return "STOP", nil
			}
			if matB == "random" {
				columnB = 0
				break
			} else {
				columnB, err = strconv.Atoi(matB)
				if err == nil {
					break
				} else {
					fmt.Println("Couldn't understand your input :(")
				}
			}

		}
		matrixB, _, _ := matrixGenerator(columnA, columnB)
		fmt.Println(matrixB)
		//Envoi des matrices au serveur
		err = sendStringMatrix("matrixA", matrixA, conn)
		if err != nil {
			return "Error", err
		}
		err = sendStringMatrix("matrixB", matrixB, conn)
		if err != nil {
			return "Error", err
		}
	}

	//Réception de la réponse du serveur
	received := false
	for !received {
		rd := bufio.NewReader(conn)
		ba, err := rd.ReadString(byte(':'))
		ba = strings.Trim(ba, ":")
		if err != nil {
			return "Error connection", err
		}

		//Vérifie que c'est la réponse du erveur
		if ba == "matrixC" {
			fmt.Println(ba)
			received = true
			//Reçoit la taille de la matrice finale
			ligne, _, err := rd.ReadRune()
			if err != nil {
				return "Error connection", err
			}
			col, _, err := rd.ReadRune()
			if err != nil {
				return "Error connection", err
			}
			fmt.Println(ligne, ":", col)
			//Réception matrice finale
			result = make([][]int, ligne)
			for i := 0; i < int(ligne); i++ {
				for j := 0; j < int(col); j++ {
					valS, err := rd.ReadString(':')
					if err != nil {
						return "Error connection", err
					}
					valT := strings.Trim(valS, ":")
					val, err := strconv.Atoi(valT)
					if err != nil {
						return "Str conversion not possible", err
					}

					result[i] = append(result[i], val)
				}
			}
			fmt.Println()
		}

	}
	//Affichage réponse
	for _, row := range result {
		for _, val := range row {
			fmt.Print(val, " ")
		}
		fmt.Println()
	}
	return "", nil

}
func sendFile(name string, f string, conn net.Conn) error {
	//Ouvre fichier
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	fmt.Println("Start sending file ", f, "!")
	wr := bufio.NewWriter(conn)
	wr.WriteString(name + ":")
	fileSize := strconv.FormatInt(fileInfo.Size(), 10)
	wr.WriteString(fileSize + ":")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wr.WriteString(scanner.Text() + "\n")

	}
	wr.WriteString("EOF\n")

	wr.Flush()
	fmt.Println("File ", f, " has been sent!")

	return nil

}
func sendStringMatrix(name string, mat string, conn net.Conn) error {

	fmt.Println("Start sending", name, "!")
	wr := bufio.NewWriter(conn)
	wr.WriteString(name + ":")
	fileSize := strconv.Itoa(len(mat) * 4)
	wr.WriteString(fileSize + ":")

	scanner := bufio.NewScanner(strings.NewReader(mat))
	for scanner.Scan() {
		wr.WriteString(scanner.Text() + "\n")
		fmt.Println(scanner.Text())
	}
	wr.WriteString("EOF\n")

	wr.Flush()
	fmt.Println(name, " has been sent!")

	return nil
}
func matrixGenerator(row int, column int) (string, int, int) {
	var matrix string = ""
	fmt.Println("row", row, "column", column)
	if row <= 0 {
		row = rand.Intn(100)
		fmt.Println("changed")
	}
	if column <= 0 {
		column = rand.Intn(100)
		fmt.Println("changed")
	}

	for i := 0; i < row; i++ {
		for j := 0; j < column-1; j++ {
			matrix += strconv.Itoa(rand.Intn(1000)) + " "
		}
		matrix += strconv.Itoa(rand.Intn(1000)) + "\n"
		//fmt.Println("aa", i, matrix)
	}
	fmt.Println("row", row, "column", column)
	return matrix, row, column
}

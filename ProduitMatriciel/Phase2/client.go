package main

import (
	"bufio"
	"fmt"
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
		fmt.Println("1. Send problem to server")
		fmt.Println("2. End connection")
		fmt.Println()
		var answer string
		fmt.Scanln(&answer)
		if answer == "1" {
			s, err := client(conn)
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

		} else if answer == "2" {
			wr := bufio.NewWriter(conn)
			wr.WriteString("STOP:")
			conn.Close()
			break
		} else {
			fmt.Println("Couldn't understand your answer. Please choose between")
		}

	}

}
func client(conn net.Conn) (string, error) {

	var result [][]int
	var matA string
	var matB string

	//Demander le fichier des matrices
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

	//Réception de la réponse du serveur
	received := false
	for !received {
		rd := bufio.NewReader(conn)
		ba, err := rd.ReadString(byte(':'))
		ba = strings.Trim(ba, ":")
		if err != nil {
			return "Error connection", err
		}
		fmt.Println(ba)

		//Vérifie que c'est la réponse du erveur
		if ba == "matrixC" {
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
					val, _, err := rd.ReadRune()
					if err != nil {
						return "Error connection", err
					}
					result[i] = append(result[i], int(val))
				}
			}
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

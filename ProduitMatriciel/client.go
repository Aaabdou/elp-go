package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func client() {
	conn, err := net.Dial("tcp", CONN_PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Print("Adress matrixA:")
		var matA string
		fmt.Scanln(&matA)
		if matA == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		err = sendFile("matrixA", matA, conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("Adress matrixB:")
		var matB string
		fmt.Scanln(&matB)
		if matB == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		err = sendFile("matrixB", matB, conn)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func sendFile(name string, f string, conn net.Conn) error {
	//Ouvre fichier
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	BUFFERSIZE := 1024

	//Envoi le nom
	conn.Write([]byte(name))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")

	//Envoi la matrice
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	file.Close()
	fmt.Println("File has been sent, closing connection!")
	return nil
}

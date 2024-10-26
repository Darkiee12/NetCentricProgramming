package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	for {
		received := make([]byte, 1024)
		n, err := conn.Read(received)
		println(string(received))
		if err != nil {
			if err == io.EOF {
				println("[Client] Connection closed by server.")
				break
			}
			println("Read data failed:", err.Error())
			os.Exit(1)
		}
		if strings.Contains(string(received[:n]), "[Response]") {
			fmt.Print("Response: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input := scanner.Text()
			_, err = conn.Write([]byte(input))
			if err != nil {
				println("Write data failed:", err.Error())
				os.Exit(1)
			}
		}
	}
}

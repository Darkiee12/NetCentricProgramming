package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8081"
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
			response(conn)
		} else if strings.Contains(string(received[:n]), "[Download]") {
			receiveFile(conn)
		}
	}
}

func receiveFile(conn net.Conn) {
	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	defer conn.Close()
	now := time.Now().UnixNano()
	file, err := os.Create(strconv.FormatInt(now, 10))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	totalBytes := 0
	buffer := make([]byte, 1024) // 1 KB buffer

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file reached")
				break
			}
			fmt.Println("Error reading from connection:", err)
			return
		}
		if n >= 6 && string(buffer[:n]) == "[Done]" {
			break
		}
		if totalBytes+n > maxFileSize {
			fmt.Println("File size limit exceeded")
			break
		}
		if _, err := file.Write(buffer[:n]); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		totalBytes += n
	}

	fmt.Println("File received successfully")
}

func response(conn net.Conn) {
	fmt.Print("Response: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	_, err := conn.Write([]byte(input))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}
	clearScreen()
}
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

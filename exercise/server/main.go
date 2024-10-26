package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()
	log.Println("Server started on", HOST+":"+PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Accepted connection from", conn.RemoteAddr())
	conn.Write([]byte("Welcome to the service\n"))
	for {
		conn.Write([]byte("[Response]Enter command (MSG_message or VLT_max_count): "))
		input, err := handleInput(conn)
		if err != nil {
			log.Println("Error reading data:", err)
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.HasPrefix(input, "MSG_") {
			message := strings.TrimPrefix(input, "MSG_")
			conn.Write([]byte("You said: " + message + "\n"))
		} else if strings.HasPrefix(input, "VLT_") {
			vltInput := strings.TrimPrefix(input, "VLT_")
			parts := strings.Split(vltInput, "_")
			if len(parts) != 2 {
				conn.Write([]byte("Error: Invalid input format for VLT command\n"))
				continue
			}
			max, err := strconv.Atoi(parts[0])
			if err != nil {
				conn.Write([]byte("Error: Invalid max number\n"))
				continue
			}
			count, err := strconv.Atoi(parts[1])
			if err != nil {
				conn.Write([]byte("Error: Invalid count number\n"))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("Your numbers are: %v\n", Vietlot(max, count))))
		} else {
			conn.Write([]byte("Error: Unknown command\n"))
		}
	}
}

func handleInput(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func Vietlot(max, count int) []int {
	if count > max {
		return []int{}
	}
	numbers := make([]int, max)
	for i := 0; i < max; i++ {
		numbers[i] = i + 1
	}
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	return numbers[:count]
}

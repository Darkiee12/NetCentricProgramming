package main

import (
	"fmt"
	"lab3/models"
	"log"
	"net"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	// Server setup and listen logic here
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			g := &models.Game{}
			g.Gameloop(c)
		}(conn)
	}
}

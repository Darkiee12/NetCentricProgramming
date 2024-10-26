package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"server/models"
	"syscall"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	// Server setup and listen logic here
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
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
			break // Exit the loop on error
		}

		go func(c net.Conn) {
			defer c.Close()
			g := &models.Game{}
			err := g.Gameloop(c)
			if err != nil {
				log.Println("Error in game loop:", err)
				return
			}
		}(conn)
	}

}

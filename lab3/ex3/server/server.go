package main

import (
	"fmt"
	"log"
	"net"
	"server/models"
)

const (
	HOST   = "localhost"
	PORT   = "8081"
	TYPE   = "tcp"
	FOLDER = "assets"
)

func main() {
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()
	fmt.Printf("Server started at %s:%s\n", HOST, PORT)
	FileManager := models.FileManager{}
	FileManager.Init()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			err := FileManager.Menu(conn, &models.Account{})
			if err != nil {
				log.Fatalf("Error handling connection: %v", err)
			}
		}(conn)

	}
}

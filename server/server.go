package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	portNumber := flag.String("port", "", "Port number")
	flag.Parse()

	if *portNumber == "" {
		fmt.Println("Usage: go run client.go --port <portNumber>")
		return
	}

	port := ":" + *portNumber
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("Server started. Listening on port %s.\n", *portNumber)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle the connection here, if needed.
	// In this example, we don't do anything with the connection.

	fmt.Println("New connection established:", conn.RemoteAddr())
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var activeConnections struct {
	sync.Mutex
	connections map[net.Conn]struct{}
}

func main() {
	portNumber := flag.String("port", "", "Port number")
	flag.Parse()

	if *portNumber == "" {
		fmt.Println("Usage: go run server.go --port <portNumber>")
		return
	}

	activeConnections.connections = make(map[net.Conn]struct{})

	port := ":" + *portNumber
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("Server started. Listening on port %s.\n", *portNumber)

	// Listen for termination signals (Ctrl-C or SIGTERM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("\nTermination signal received. Closing connections...")
		activeConnections.Lock()
		for conn := range activeConnections.connections {
			conn.Close()
		}
		activeConnections.Unlock()
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		activeConnections.Lock()
		activeConnections.connections[conn] = struct{}{}
		activeConnections.Unlock()

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Do whatever operations you need to perform with the connection.
	// In this example, we keep the connection open until termination.

	fmt.Println("New connection established:", conn.RemoteAddr())

	// Wait for termination signal (SIGTERM) before closing the connection.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs

	conn.Close()

	activeConnections.Lock()
	delete(activeConnections.connections, conn)
	activeConnections.Unlock()

	fmt.Println("Connection closed:", conn.RemoteAddr())
}

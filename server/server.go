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
	xchng := flag.String("xchng", "", "Data exchange to use (0 = no exchange, 1 = client send only, 2 = client/server response")
	flag.Parse()

	if *portNumber == "" || *xchng == "" {
		fmt.Println("Usage: go run server.go --port <portNumber> --xchng <0|1|2>")
		return
	}

	activeConnections.connections = make(map[net.Conn]struct{})

	port := ":" + *portNumber
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("Server started. Listening on port %s. Message exchange mode is %s\n", *portNumber, *xchng)

	// Listen for termination signals (Ctrl-C or SIGTERM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Create a thread to listen for termination signal and then gracefully close the connections.
	go func() {
		<-signalCh
		fmt.Println("\nTermination signal received. Closing connections...")
		activeConnections.Lock()
		for conn := range activeConnections.connections {
			log.Printf("Closing connection: %s", conn.RemoteAddr())
			conn.Close()
		}
		activeConnections.Unlock()
		listener.Close()
		os.Exit(0)
	}()

	// Loop forever and accept any incoming connection.
	// Each connection will have its own handler.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		activeConnections.Lock()
		activeConnections.connections[conn] = struct{}{}
		activeConnections.Unlock()

		// Start a thread to handle the connection
		go handleConnection(conn, xchng)
	}
}

func handleConnection(conn net.Conn, xchng *string) {
	// Do whatever operations you need to perform with the connection.
	// In this example, we keep the connection open until termination.

	fmt.Println("New connection established:", conn.RemoteAddr())

	// Loop forever and respond to hello messages
	for {
		if *xchng == "1" || *xchng == "2" {
			buf := make([]byte, 5)
			_, err := conn.Read(buf)
			if err != nil {
				log.Printf("Failed to read from connection: %v", err)
				return
			}

			request := string(buf)
			if request == "hello" {
				log.Printf("Hello from: %s", conn.RemoteAddr())

				// Check if we are supposed to send a response back.
				if *xchng == "2" {
					_, err = conn.Write([]byte("world"))
					if err != nil {
						log.Printf("Failed to send response: %v", err)
					}

					log.Printf("Sent 'world' response to client: %s", conn.RemoteAddr())
				}
			} else {
				log.Printf("Received an unexpected request (%s) from client: %s", request, conn.RemoteAddr())
			}
		} else {
			select {} // Enter an idle state and consume minimal resources
			log.Printf("ERROR: Idle state has reached unreachable code!!!!!")

		}
	}
}

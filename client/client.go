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

var wg sync.WaitGroup
var connections []net.Conn

func cleanup() {
	fmt.Printf("Program exiting. Cleaning up conections...")

	fmt.Println("\nTermination signal received. Closing connections...")
	for _, conn := range connections {
		log.Printf("Closing connection: %s", conn.RemoteAddr())
		conn.Close()
	}
	os.Exit(0)
}

func main() {
	defer cleanup()

	numConnections := flag.Int("numcons", 0, "Number of connections to create")
	portNumber := flag.String("port", "", "Port number")
	url := flag.String("url", "", "URL to use")
	xchng := flag.String("xchng", "", "Data exchange to use (0 = no exchange, 1 = client send only, 2 = client/server response")
	flag.Parse()

	if *numConnections <= 0 || *portNumber == "" || *url == "" || *xchng == "" {
		log.Printf("Usage: go run client.go --numcons <numConnections> --port <portNumber> --url <url> --xchng <0|1|2>")
		return
	}

	fmt.Printf("Starting client. Num connections: %d, URL: %s, Port: %s\n", *numConnections, *url, *portNumber)

	// Listen for termination signals (Ctrl-C)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Create a thread to listen for termination signal and then gracefully close the connections.
	go func() {
		<-signalCh
		fmt.Println("\nTermination signal received. Closing connections...")
		for _, conn := range connections {
			log.Printf("Closing connection: %s", conn.RemoteAddr())
			conn.Close()
		}
		os.Exit(0)
	}()

	for i := 0; i < *numConnections; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *url, *portNumber))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			continue
		}
		defer conn.Close()

		connections = append(connections, conn)
		fmt.Println("Connection", i+1, "established:", conn.RemoteAddr())

		if *xchng == "1" || *xchng == "2" {
			// Send "hello" message to the server
			_, err = conn.Write([]byte("hello"))
			if err != nil {
				log.Fatalf("Failed to send message: %v", err)
			}
			log.Printf("Sent 'hello' message to server: %s", conn.RemoteAddr())

			if *xchng == "2" {
				go handleResponse(conn)
			}
		}
	}

	// Sit idle and keep connections open but don't consume any more resources
	select {}
}

func handleResponse(conn net.Conn) {
	buf := make([]byte, 5)
	log.Printf("Waiting for response from %s", conn.RemoteAddr())

	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read from connection: %v", err)
		return
	}

	resp := string(buf)
	if resp == "world" {
		log.Printf("Received response (%s) from server: %s", resp, conn.RemoteAddr())
	} else {
		log.Printf("ERROR: Received an unexpected response (%s) from server: %s", resp, conn.RemoteAddr())
	}
}

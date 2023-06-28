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

func main() {
	numConnections := flag.Int("numcons", 0, "Number of connections to create")
	portNumber := flag.String("port", "", "Port number")
	url := flag.String("url", "", "URL to use")
	flag.Parse()

	if *numConnections <= 0 || *portNumber == "" || *url == "" {
		fmt.Println("Usage: go run client.go --numcons <numConnections> --port <portNumber> --url <url>")
		return
	}

	fmt.Printf("Starting client. Num connections: %d, URL: %s, Port: %s\n", *numConnections, *url, *portNumber)

	var connections []net.Conn
	var wg sync.WaitGroup

	// Listen for termination signals (Ctrl-C)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("\nTermination signal received. Closing connections...")
		for _, conn := range connections {
			conn.Close()
		}
		wg.Wait()
		os.Exit(0)
	}()

	for i := 0; i < *numConnections; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *url, *portNumber))
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			continue
		}

		connections = append(connections, conn)
		fmt.Println("Connection", i+1, "established:", conn.RemoteAddr())

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()

			// Keep the connection open until termination
			<-signalCh
			fmt.Println("Closing connection:", conn.RemoteAddr())
			conn.Close()
		}(conn)
	}

	wg.Wait()
	fmt.Println("All connections closed. Exiting.")
}

package main

import (
	"net"
	"netcat/funcs"
	"os"
	"testing"
	"time"
)

// test for both server startup and client connection
func TestTCPServer(t *testing.T) {
	os.Args = []string{".", "8080"}
	go main()                   // Start the server in a separate goroutine
	time.Sleep(1 * time.Second) // Wait for the server to start

	// Try to connect a client to the server
	conn, err := net.Dial("tcp", funcs.GetLocalIP()+":8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// If we successfully connected, the server is up and running
	t.Log("Server is up and running")
}

package main

import (
	"fmt"
	"log"
	"net"
	"netcat/funcs"
	"os"
)

func main() {
	PORT := 8989
	TYPE := "tcp"

	if len(os.Args) == 2 {
		PORT, _ = funcs.Atoi(os.Args[1])
		_, err := funcs.Atoi(os.Args[1])

		if err != nil {
			log.Fatalf("[USAGE]: ./TCPchat $PORT\n")
		}

	} else if len(os.Args) > 2 {
		log.Fatalf("[USAGE]: ./TCPchat $PORT\n")
	}

	listenAddr := fmt.Sprintf("%s:%d", funcs.GetLocalIP(), PORT)

	listen, err := net.Listen(TYPE, listenAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("Listening on %s...\n", listenAddr)
	defer listen.Close()

	go funcs.HandleShutdown()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go funcs.HandleConnection(conn)
	}
}

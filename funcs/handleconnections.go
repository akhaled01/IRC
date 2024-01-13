package funcs

import (
	"net"
)

func HandleConnection(conn net.Conn) {
	// limit connections to 10 connections only
	if len(userpool) == 10 {
		conn.Write([]byte("User limit reached, try again later\n"))
		conn.Close()
		return
	}
	defer conn.Close()

	// Create a unique user ID for this connection
	userID := genID()
	user := User{ID: userID, Connection: conn}

	// Prompt the user for their name and perform authentication
	if !authenticateUser(&user) {
		return
	}


	// Add the user to the user pool
	mu.Lock()
	userpool = append(userpool, user)
	mu.Unlock()

	// Handle user messages
	messageCh := make(chan string, 1)
	go readUserInput(&user, messageCh)

	for {
		select {
		case message := <-messageCh:
			// Broadcast the message to all users
			broadcastMessage(user, message)
		case <-shutdownCh:
			return
		}
	}
}

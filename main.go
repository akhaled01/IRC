package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"netcat/funcs"
	"os"
	"os/signal"
	"sync"
	"time"
)

type User struct {
	ID         string
	Name       string
	Connection net.Conn
}

var (
	mu         sync.Mutex
	userpool   []User
	shutdownCh chan struct{} // Channel for graceful shutdown
	useridseq  int
	textlog    string
	file, _    = os.Create("textlog.txt")
)

func genID() string {
	useridseq++
	return fmt.Sprintf("user%d", useridseq)
}

func main() {
	PORT := 8989
	TYPE := "tcp"

	if len(os.Args) > 1 {
		PORT, _ = funcs.Atoi(os.Args[1])
		_, err := funcs.Atoi(os.Args[1])

		if err != nil {
			log.Fatalf("[USAGE]: ./server.go $PORT\n")
		}

	}

	listenAddr := fmt.Sprintf("%s:%d", funcs.GetLocalIP(), PORT)
	fmt.Printf("Listening on %s...\n", listenAddr)

	listen, err := net.Listen(TYPE, listenAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listen.Close()

	shutdownCh = make(chan struct{})
	go handleShutdown()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Block until a signal is received
	sig := <-sigCh
	log.Printf("Received signal %s. Shutting down...\n", sig)
	close(shutdownCh) // Signal other goroutines to stop
	// Perform any cleanup here.

	os.Exit(0)
}

func writetextlog(conn net.Conn) {
	conn.Write([]byte(textlog))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a unique user ID for this connection
	userID := genID()
	user := User{ID: userID, Connection: conn}

	// Prompt the user for their name and perform authentication
	if !authenticateUser(&user) {
		return
	}

	// export text log

	writetextlog(user.Connection)

	// Add the user to the user pool
	mu.Lock()
	userpool = append(userpool, user)
	mu.Unlock()

	defer func() {
		// Remove the user from the user pool and broadcast their departure
		mu.Lock()
		removeUser(user)
		mu.Unlock()
	}()

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

func authenticateUser(user *User) bool {
	timeoutDuration := 30 * time.Second
	conn := user.Connection

	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	defer conn.SetReadDeadline(time.Time{})

	logo := `
	 .--.
	|o_o |
	|:_/ |
   //   \ \
  (|     | )
 /'\_   _/ \
 \___)=(___0\/
`
	conn.Write([]byte(logo))
	conn.Write([]byte("\nWelcome to net majles\n"))
	conn.Write([]byte("Enter Your name: "))

	bufReader := bufio.NewReader(conn)
	bytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		log.Println("Error reading user input:", err)
		return false
	}

	user.Name = string(bytes[:len(bytes)-1])
	if user.Name == "" {
		log.Println("Empty username. Closing connection.")
		return false
	}

	// Broadcast user's arrival
	broadcastMessage(*user, fmt.Sprintf("%s joined the chat\n", user.Name))

	return true
}

func readUserInput(user *User, messageCh chan<- string) {
	conn := user.Connection
	bufReader := bufio.NewReader(conn)

	for {
		message, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				removeUser(*user)
				return
			} else {
				fmt.Printf("Error reading input %s\n", err)
				return
			}
		}
		messageCh <- message
	}
}

func broadcastMessage(sender User, message string) {
	t := time.Now().Format(time.ANSIC)
	message = fmt.Sprintf("[%v][%v] %v", t, sender.Name, message)
	file.Write([]byte(message))
	textlog += message
	mu.Lock()
	for _, u := range userpool {
		if u.ID != sender.ID {
			u.Connection.Write([]byte(message))
		}
	}
	mu.Unlock()
}

func removeUser(user User) {
	// Broadcast user's departure
	broadcastMessage(user, fmt.Sprintf("%s left the chat\n", user.Name))
	for i, u := range userpool {
		if u.ID == user.ID {
			// Close the user's connection
			user.Connection.Close()

			// Remove the user from the user pool
			userpool = append(userpool[:i], userpool[i+1:]...)

			return
		}
	}
}

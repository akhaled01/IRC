package funcs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func authenticateUser(user *User) bool {
	// Set a timeout duration for user authentication
	timeoutDuration := 30 * time.Second
	conn := user.Connection

	// Set a read deadline to enforce the timeout
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	// Display a welcome message and logo to the user
	logo := `
     .--.
    |o_o |
    |:_/ |
   //   \ \
  (| EAA | )
 /'\_   _/ \
 \___)=(___0\/
`
	conn.Write([]byte(logo))
	conn.Write([]byte("\nWelcome to net majles\n"))
	conn.Write([]byte("Enter Your name: "))

	// Create a buffered reader to read the user's input
	bufReader := bufio.NewReader(conn)
	bytes, err := bufReader.ReadBytes('\n')

	// Check for errors while reading user input
	if err == io.EOF {
		return false
	} else if err, ok := err.(net.Error); ok && err.Timeout() {
		user.Connection.Write([]byte("\ntime runout try agin leter\n"))
		return false
	} else if err != nil {
		log.Println("Error reading user input:", err)
		return false
	}

	// Remove the trailing newline character and set the user's name
	user.Name = string(bytes[:len(bytes)-1])

	// Check if the username is empty
	if strings.TrimSpace(user.Name) == "" {
		log.Println("Empty username. Closing connection.")
		user.Connection.Write([]byte("Please Input a username\n"))
		return false
	}
	conn.SetReadDeadline(time.Time{})
	// export text log
	writetextlog(user.Connection)
	conn.Write([]byte("\n----------- HISTORY -----------\n"))
	// Broadcast user's arrival to the chat
	broadcastMessage(*user, fmt.Sprintf("%s joined the chat\n", user.Name))

	return true
}

func genID() string {
	useridseq++
	return fmt.Sprintf("user%d", useridseq)
}

// when the user enters /un
func changeUsername(u *User) {
	conn := u.Connection

	// Set a read deadline for the connection
	timeoutDuration := 30 * time.Second
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	defer conn.SetReadDeadline(time.Time{})

	conn.Write([]byte("Enter Your new name: "))

	// Read user input from the connection
	bufReader := bufio.NewReader(conn)
	bytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		log.Println("Error reading user input:", err)
		return
	}

	// Update the user's name with the input
	newname := string(bytes[:len(bytes)-1])

	// Check if the username is empty or contains only whitespace the program will keep the current username
	if strings.TrimSpace(newname) == "" {
		u.Connection.Write([]byte("Invalid Username\n"))
		return
	} else {
		u.Name = newname
	}
}

func removeUser(user User) {
	// Broadcast user's departure
	broadcastMessage(user, fmt.Sprintf("%s left the chat\n", user.Name))
	mu.Lock()
	for i, u := range userpool {
		if u.ID == user.ID {
			// Close the user's connection
			user.Connection.Close()

			// Remove the user from the user pool
			userpool = append(userpool[:i], userpool[i+1:]...)
			mu.Unlock()

			return
		}
	}
}
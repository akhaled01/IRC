package funcs

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func readUserInput(user *User, messageCh chan<- string) {
    // Get the user's connection and create a buffered reader for reading input
    conn := user.Connection
    bufReader := bufio.NewReader(conn)

    for {
        // Read a message from the user
        message, err := bufReader.ReadString('\n')

        // Check for errors while reading
        if err != nil {
            // If it's an EOF error, remove the user and return, indicating they left the chat
            if err == io.EOF {
                removeUser(*user)
                return
            } else {
                fmt.Printf("Error reading input %s\n", err)
                return
            }
        }

        // Check if the message exceeds the word limit (255 characters)
        if len(message) > 255 {
            user.Connection.Write([]byte("Exceeded word limit (255 characters only)\n"))
            continue
        }

        // Check if the message is a special command
        if strings.TrimSpace(message) == "/nu" {
            oldUsername := user.Name
            changeUsername(user)

            // If the username was changed, notify the chat
            if user.Name != oldUsername {
                message = fmt.Sprintf("%s changed his username to %s\n", oldUsername, user.Name)
                messageCh <- message
            }
        // } else if strings.TrimSpace(message) == "/l" {
        //     // User wants to leave the chat, remove them
        //     removeUser(*user)
		// 	return
        } else if strings.TrimSpace(message) == "/h" {
            // User requested help, provide a list of implemented commands
            helpMessage := `
            Implemented commands:
            /l : leave the chat
            /nu : change username
            `
            user.Connection.Write([]byte(helpMessage))
        } else if strings.TrimSpace(message) != "" {
            // If it's not a special command and not empty, send it to the chat
            messageCh <- message
        }
    }
}

func broadcastMessage(sender User, message string) {
	// Get current time and format it as a string and Format the message with current time, sender's name, and the message
	t := time.Now().Format(time.ANSIC) 
	message = fmt.Sprintf("[%v][%v] %v", t, sender.Name, message)

	// Write the formatted message to a file and append it to textlog variable
	file.Write([]byte(message))  
	textlog += message 

	// Write the message to the to each user but not the sender
	mu.Lock()  
	for _, u := range userpool { 
		if u.ID != sender.ID { 
			u.Connection.Write([]byte(message)) 
		}
	}
	mu.Unlock()  
}

// display the history of the messages to the new user
func writetextlog(conn net.Conn) {
	conn.Write([]byte(textlog))
}

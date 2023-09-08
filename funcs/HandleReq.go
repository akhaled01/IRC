package funcs

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type User struct {
	name       string
	connection net.Conn
}

var mu sync.Mutex
var userpool []User

func Reply(reply []byte) {
	for _, c := range userpool {
		c.connection.Write(reply)
	}
}

func AddUser(listen net.Listener) {
	// lock the thread to accept connections, then unlock
	mu.Lock()
	conn, err := listen.Accept()
	AuthenticateUser(conn)
	mu.Unlock()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//TODO: Add username authentication here
	channel := make(chan string, 1)
	// for _, u := range userpool {
	// 	go HandleRequest(conn, channel, u)
	// }

	// go HandleRequest(conn, channel, u)
	// for {
	// incoming request
	buffer := make([]byte, 1024)
	_, err1 := conn.Read(buffer)
	if err1 != nil {
		channel <- fmt.Sprintf("\n left the chat 11\n")

	}
	responseStr := ""
	// write data to channel

	for ur := range userpool {
		t := time.Now().Format(time.ANSIC)
		responseStr = fmt.Sprintf("\n[%v][%v] %v\n", t, ur, string(buffer[:]))
	}

	// t := time.Now().Format(time.ANSIC)
	// responseStr := fmt.Sprintf("\n[%v] %v\n", t, string(buffer[:]))

	channel <- responseStr
	for message := range channel {
		fmt.Println(message)
		// mu.Lock()
		go Reply([]byte(message))
		// mu.Unlock()
	}
	// }

}

func AuthenticateUser(conn net.Conn) {
	timeoutDuration := 30 * time.Second
	bufReader := bufio.NewReader(conn)
	// Set a deadline for reading. Read operation will fail if no data
	// is received after deadline.
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
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
	// Read tokens delimited by newline
	bytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("\nNo username detected. Exiting....\n"))
		conn.Close()
		return
	}
	conn.SetDeadline(time.Time{})
	newUser := &User{}
	newUser.name = string(bytes[:len(bytes)-1])
	if newUser.name == "" {
		conn.Close()
	}
	newUser.connection = conn
	Reply([]byte(fmt.Sprintf("\n%v joined the chat\n", newUser.name)))
	// append to users
	userpool = append(userpool, *newUser)

}

func HandleRequest(conn net.Conn, channel chan string, u User) {
	for {
		// incoming request
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			CloseConnection(u)
			break
		}

		// write data to channel
		t := time.Now().Format(time.ANSIC)
		responseStr := fmt.Sprintf("\n[%v][%v] %v\n", t, u.name, string(buffer[:]))

		channel <- responseStr
	}
}

func CloseConnection(toremoveconn User) {
	mu.Lock()
	for index, v := range userpool {
		if v == toremoveconn {
			toremoveconn.connection.Close()
			Reply([]byte(fmt.Sprintf("\n%v left the chat\n", userpool[index].name)))
			userpool = append(userpool[:index], userpool[index+1:]...)
			break
		}
	}
	mu.Unlock()

	fmt.Println("After ", userpool)
	fmt.Println("............................................................")
	fmt.Println("Connection closed and removed successfully")
}

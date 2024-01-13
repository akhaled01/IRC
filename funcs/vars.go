package funcs

import (
	"net"
	"os"
	"sync"
)

type User struct {
	ID         string
	Name       string
	Connection net.Conn
}

var (
	mu         sync.Mutex
	userpool   []User
	shutdownCh = make(chan struct{}) // Channel for graceful shutdown
	useridseq  int
	textlog    string
	file, _    = os.Create("textlog.txt")
)

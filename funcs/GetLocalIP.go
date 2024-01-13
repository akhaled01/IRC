package funcs

import (
	"log"
	"net"
)

// GetLocalIP retrieves the local IP address by establishing a UDP connection to a known IP address (e.g., Google's DNS server).
// It returns the local IP address as a string.
func GetLocalIP() string {
	// Establish a UDP connection to a known IP address (e.g., Google's DNS server)
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Retrieve the local address information from the connection
	localAddress := conn.LocalAddr().(*net.UDPAddr)

	// Return the local IP address as a string
	return localAddress.IP.String()
}

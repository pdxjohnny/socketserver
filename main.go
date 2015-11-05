package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/pdxjohnny/socketserver/server"
)

func connection(conn net.Conn) {
	// will listen for message to process ending in newline (\n)
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading from connection", err)
		return
	}
	// output message received
	fmt.Print("Message Received:", string(message))
	// sample process for string received
	newmessage := strings.ToUpper(message)
	// send new string back to client
	conn.Write([]byte(newmessage + "\n"))
}

func main() {
	server.StartServer(connection)
}

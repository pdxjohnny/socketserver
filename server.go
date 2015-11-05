package main

import (
	"errors"
	"log"
	"net"

	"github.com/pdxjohnny/socketserver/server"
)

const (
	NumRequests   = 10000
	CreateClients = 1000
)

var (
	ClientPool chan chan bool
)

type TestClient struct {
	clientPool chan chan bool
	doReq      chan bool
	Conn       net.Conn
}

func NewTestClient(clientPool chan chan bool) *TestClient {
	return &TestClient{
		clientPool: clientPool,
		doReq:      make(chan bool),
		Conn:       nil,
	}
}

func (c *TestClient) Run() {
	go func() {
		for {
			// log.Println("TestClient running")
			c.clientPool <- c.doReq

			// log.Println("TestClient added to pool")
			select {
			case makeReq := <-c.doReq:
				if !makeReq {
					panic(errors.New("TestClient was asked to exit"))
				}
				// log.Println("Make a connection")
				conn, err := net.Dial("tcp", "127.0.0.1:25001")
				if err != nil {
					panic(err)
				}
				// log.Println("Connection established")
				conn.Close()
			}
		}
	}()
}

func main() {
	done := make(chan bool)
	counter := 0
	go server.StartServer(func(conn net.Conn) {
		conn.Write([]byte("Hello\n"))
		conn.Close()
		counter++
		log.Println("Sent back hello", counter)
		if counter >= NumRequests {
			done <- true
		}
	})
	ClientPool = make(chan chan bool)
	for index := 0; index < CreateClients; index++ {
		log.Println("Creating client", index)
		client := NewTestClient(ClientPool)
		client.Run()
	}
	log.Println("Server started")
	for index := 0; index < NumRequests; index++ {
		doReq := <-ClientPool
		doReq <- true
	}
	log.Println("Waiting for requests to finish")
	<-done
	log.Println("All requests finished")
}

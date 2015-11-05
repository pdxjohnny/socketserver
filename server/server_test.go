package server

import (
	"log"
	"net"
	"testing"
)

const (
	NumRequests   = 1
	CreateClients = 1
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
			c.clientPool <- c.doReq

			select {
			case makeReq := <-c.doReq:
				if !makeReq {
					log.Println("TestClient Exiting")
					return
				}
				if c.Conn == nil {
					conn, err := net.Dial("tcp", "127.0.0.1:25001")
					if err != nil {
						log.Println(err)
						return
					}
					c.Conn = conn
				}
				log.Println("Make a request")
			}
		}
	}()
}

func TestLoad(t *testing.T) {
	for index := 0; index < CreateClients; index++ {
		log.Println("Creating client", index)
		NewTestClient(ClientPool)
	}
	StartServer(func(conn net.Conn) {
		conn.Write([]byte("Hello\n"))
	})
	for index := 0; index < NumRequests; index++ {
		log.Println("Sending req")
		doReq := <-ClientPool
		doReq <- true
	}
}

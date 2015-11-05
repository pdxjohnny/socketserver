package server

import (
	"net"
	"os"
	"strconv"
)

const (
	// DefaultMaxWorkers used if MaxWorker is empty
	DefaultMaxWorkers = 8
)

var (
	// MaxWorker -
	MaxWorker, _ = strconv.Atoi(os.Getenv("MAX_WORKERS"))
	// MaxQueue -
	MaxQueue, _ = strconv.Atoi(os.Getenv("MAX_QUEUE"))
)

// Dispatcher -
type Dispatcher struct {
	// ConnQueue -
	// A buffered channel that we can send work requests on.
	ConnQueue chan net.Conn
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan net.Conn
	// The max number of workers we can create and dispatch to
	MaxWorkers int
	// Handler - fucntions workers call
	Handler func(conn net.Conn)
}

// NewDispatcher -
func NewDispatcher(maxWorkers int, handler func(conn net.Conn)) *Dispatcher {
	return &Dispatcher{
		ConnQueue:  make(chan net.Conn),
		WorkerPool: make(chan chan net.Conn, maxWorkers),
		MaxWorkers: maxWorkers,
		Handler:    handler,
	}
}

// Run -
func (d *Dispatcher) Run() {
	// log.Println("Starting", d.MaxWorkers, "workers")
	// Starting MaxWorkers number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.Handler)
		worker.Start()
	}

	// log.Println("Started", d.MaxWorkers, "workers")
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		// log.Println("Waiting for conn from ConnQueue")
		select {
		case conn := <-d.ConnQueue:
			// log.Println("Got conn from ConnQueue")
			// a conn request has been received
			go func(conn net.Conn) {
				// try to obtain a worker conn channel that is available.
				// this will block until a worker is idle
				connChannel := <-d.WorkerPool

				// dispatch the conn to the worker conn channel
				connChannel <- conn
			}(conn)
		}
	}
}

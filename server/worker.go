package server

import "net"

// Worker represents the worker that executes the conn
type Worker struct {
	WorkerPool  chan chan net.Conn
	ConnChannel chan net.Conn
	quit        chan bool
	// Handler - fucntion to call
	Handler func(conn net.Conn)
}

// NewWorker -
func NewWorker(workerPool chan chan net.Conn, handler func(conn net.Conn)) Worker {
	return Worker{
		WorkerPool:  workerPool,
		ConnChannel: make(chan net.Conn),
		quit:        make(chan bool),
		Handler:     handler,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.ConnChannel

			select {
			case conn := <-w.ConnChannel:
				// log.Println("Got a conn", conn)
				if w.Handler != nil {
					// log.Println("Calling handle")
					w.Handler(conn)
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

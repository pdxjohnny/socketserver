package server

import "net"

func startDispatcher(handler func(conn net.Conn)) *Dispatcher {
	// log.Println("MaxWorker is", MaxWorker)
	if MaxWorker == 0 {
		MaxWorker = DefaultMaxWorkers
		// log.Println("MaxWorker is", MaxWorker)
	}
	dispatcher := NewDispatcher(MaxWorker, handler)
	dispatcher.Run()
	return dispatcher
}

// StartServer -
func StartServer(handler func(conn net.Conn)) {
	// fmt.Println("Creating dispatcher...")
	dispatcher := startDispatcher(handler)
	// fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":25001")
	if err != nil {
		return
	}
	// run loop forever (or until ctrl-c)
	for {
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			// log.Println("While accepting", err)
		}
		// log.Println("Sending conn to ConnQueue")
		dispatcher.ConnQueue <- conn
	}
}

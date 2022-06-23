package main

import (
	"net"
	"net/rpc"
	"strconv"

	"github.com/kaiya/goutils/tinyrpc/core"
)

func main() {
	// Publish our Handler methods
	rpc.Register(&core.Handler{})
	Port := 3344

	// Create a TCP listener that will listen on `Port`
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(Port))

	// Close the listener whenever we stop
	defer listener.Close()

	// Wait for incoming connections
	rpc.Accept(listener)
}

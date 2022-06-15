package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kaiya/goutils/tinyrpc/rpc"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Make a buffer to hold incoming data.
	/*
		totalsizebuf := make([]byte, 4)
		_, err := io.ReadFull(conn, totalsizebuf)
		if err != nil {
			fmt.Println("Error reading totalsize:", err.Error())
			return
		}
		totalsize := binary.BigEndian.Uint32(totalsizebuf)
		fmt.Println("totalsize: ", totalsize)

		magicbuf := make([]byte, 4)
		_, err = io.ReadFull(conn, magicbuf)
		if err != nil {
			fmt.Println("Error reading magic:", err.Error())
			return
		}
		fmt.Println("magic: ", string(magicbuf))
	*/

	// Send a response back to person contacting us.
	//conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.

	msg, err := rpc.DecodePacket(conn)
	if err != nil {
		fmt.Printf("decode packet error:%s\n", err)
		return
	}

	fmt.Println("got msg len: ", len(msg))
	fmt.Println("msg: ", string(msg))
}

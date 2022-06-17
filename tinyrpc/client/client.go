package main

import (
	"fmt"
	"net/rpc"
	"strconv"

	"github.com/Kaiya/goutils/tinyrpc/core"
)

var (
	Port     = 3344
	addr     = "127.0.0.1:" + strconv.Itoa(Port)
	request  = &core.Request{Name: "world"}
	response = new(core.Response)
)

func main() {

	// Establish the connection to the adddress of the
	// RPC server
	client, _ := rpc.Dial("tcp", addr)
	defer client.Close()

	// Perform a procedure call (core.HandlerName == Handler.Execute)
	// with the Request as specified and a pointer to a response
	// to have our response back.
	err := client.Call(core.HandlerName, request, response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.Message)
}

package rpc

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"net"
	"testing"
)

func Test_Encode(t *testing.T) {
	conn, _ := net.Dial("tcp", "127.0.0.1:3333")
	defer conn.Close()

	zip := zlib.NewWriter(conn)
	bufconn := bufio.NewWriter(zip)
	err := EncodePacket(bufconn, []byte("hello"))
	if err != nil {
		fmt.Printf("encode error:%s", err)
	}
	bufconn.Flush()
	zip.Flush()
}

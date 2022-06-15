package rpc

import (
	"bufio"
	"net"
	"testing"
)

func Test_Encode(t *testing.T) {
	conn, _ := net.Dial("tcp", "127.0.0.1:3333")
	//zlib.NewWriter(conn)
	bufconn := bufio.NewWriter(conn)
	err := EncodePacket(bufconn, []byte("hello"))
	if err != nil {
		t.Errorf("encode error:%s", err)
	}
}

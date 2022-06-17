package rpc

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"net"
	"sync"
	"testing"
)

func Test_Encode(t *testing.T) {
	testStr := "hello"
	var wg sync.WaitGroup
	server, client := net.Pipe()
	// server side
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer server.Close()
		zipR, err := zlib.NewReader(server)
		if err != nil {
			t.Errorf("init zlib reader failed:%s", err)
		}
		bufconn := bufio.NewReader(zipR)
		msg, err := DecodePacket(bufconn)
		if err != nil {
			fmt.Printf("decode packet error:%s\n", err)
			return
		}
		fmt.Printf("got msg len:%d, %s ", len(msg), string(msg))
		if res := string(msg); res != testStr {
			t.Errorf("res:%s, want:%s", res, testStr)
		}
	}()

	// client side
	zip := zlib.NewWriter(client)
	bufconn := bufio.NewWriter(zip)
	err := EncodePacket(bufconn, []byte(testStr))
	if err != nil {
		fmt.Printf("encode error:%s", err)
	}
	bufconn.Flush()
	zip.Flush()
	client.Close()
	wg.Wait()
}

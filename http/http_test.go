package http

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"gitlab.momoso.com/cm/kit/pkg/basic/logrender"
	"gitlab.momoso.com/cm/kit/pkg/net/httpclient"
)

func TestTinyHttp(t *testing.T) {
	go func() {
		//server
		lis, err := net.Listen("tcp4", "127.0.0.1:3000")
		if err != nil {
			t.Errorf("listen error:%s", err)
		}
		// for {
		conn, err := lis.Accept()
		if err != nil {
			t.Errorf("accept error:%s", err)
			// break
		}
		/*
			buf := make([]byte, 10240)
			for {
				n, err := conn.Read(buf)
				t.Logf("got len:%d, err:%v, str:%s", n, err, string(buf))

				if err == io.EOF {
					break
				}
			}
		*/

		// read using io.readall
		reqBytes, err := io.ReadAll(conn)
		if err != nil {
			t.Errorf("readall from req error:%s", err)

		}
		t.Logf("req:%s", string(reqBytes))
		// read from scanner
		/*
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				t.Logf("got from conn:%s", scanner.Text())
			}
		*/

		// t.Logf("scanner error:%s", scanner.Err())
		/*
			reader := bufio.NewReader(conn)

			bytes, err := reader.ReadBytes('\n')

			if err != nil {
				t.Errorf("read bytes error:%s", err)
			}
			t.Logf("got from conn:%s", string(bytes))
		*/
		// }
	}()

	//client
	/*
		_, err := http.Get("http://127.0.0.1:3000/test")
		if err != nil {
			t.Errorf("http get error:%s", err)
		}
	*/
	httpclient.NewClient(&httpclient.Config{
		RequestTimeout: 10 * time.Second,
		Render: &logrender.Config{
			Stdout:        true,
			StdoutPattern: "*",
		},
	}).PostJSON(context.Background(), "http://127.0.0.1:3000/post-json", nil, nil, map[string]string{"json-key": "json-value"})
	time.Sleep(30 * time.Second)

}

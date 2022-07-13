package http

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/url"
)

type Response struct {
	conn       *net.Conn
	Req        *Request
	Header     map[string][]string
	StatusCode int
	reqBody    io.ReadCloser
	// Body
	w *bufio.Writer
}

type Request struct {
	Proto   string
	Url     *url.URL
	Header  map[string][]string
	Method  string
	Payload []byte
	Body    io.ReadCloser
	Host    string
	Ctx     context.Context
}

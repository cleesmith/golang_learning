package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type DelayedStartReader struct {
	Delay     time.Duration
	Content   string
	curOffset int
}

// Read will return some bytes from the Content string into the given slice. If this is the first
// time Read is called, it will first sleep for the requested delay period.
func (t *DelayedStartReader) Read(p []byte) (n int, err error) {
	if t.curOffset == 0 {
		time.Sleep(t.Delay)
	} else if t.curOffset >= len(t.Content) {
		return 0, io.EOF
	}

	bytesCopied := copy(p, t.Content[t.curOffset:])
	t.curOffset += bytesCopied
	return bytesCopied, nil
}

func MakeDelayedStartReader(delay time.Duration, content string) io.Reader {
	return &DelayedStartReader{Delay: delay, Content: content}
}

// MakeResponse constructs a *http.Response based on some input.
func MakeResponse(req *http.Request, status int, firstLine string, body string) *http.Response {
	return &http.Response{
		Status:        firstLine,
		StatusCode:    status,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          ioutil.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
	}

}

// FixHttp10Response maybe downgrades an http.Response object to HTTP 1.0, which is necessary in
// the case where the original client request was 1.0.
func FixHttp10Response(resp *http.Response, req *http.Request) {
	if req.ProtoMajor == 1 && req.ProtoMinor == 1 {
		return
	}

	resp.Proto = "HTTP/1.0"
	resp.ProtoMajor = 1
	resp.ProtoMinor = 0

	if strings.Contains(strings.ToLower(req.Header.Get("Connection")), "keep-alive") {
		resp.Header.Set("Connection", "keep-alive")
		resp.Close = false
	} else {
		resp.Close = true
	}
}

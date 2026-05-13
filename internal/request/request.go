package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestState int

const (
	Request_Initialized RequestState = iota
	Request_Done
)
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, bufferSize)
	request := &Request{
		RequestState: Request_Initialized,
	}
	readToIndex := 0

	for request.RequestState != Request_Initialized {

		if len(buffer) == cap(buffer) {
			newBuffer := make([]byte, bufferSize*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		r, err := reader.Read(buffer[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error reading from buffer")
		}
		if errors.Is(err, io.EOF) {
			request.RequestState = Request_Done
		}
		readToIndex = r

		bytesRead, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing bytes: %s Total bytes read: %d", err, bytesRead)
		}

		copy(buffer, make([]byte, 8))
		readToIndex -= bytesRead

	}

	return request, nil
}

func parseRequestLine(lines []byte) (int, *RequestLine, error) {

	idx := bytes.Index(lines, []byte("\r\n"))
	if idx == -1 {
		return 0, nil, nil
	}

	req := string(lines[:idx])
	r := strings.Split(req, " ")
	if len(r) < 3 {
		fmt.Println(r)
		return 0, nil, fmt.Errorf("Invalid HTTP request-line: %s", r)
	}

	method := r[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return 0, nil, fmt.Errorf("error: Invalid method in request-line")
		}
	}

	target := r[1]
	version := strings.TrimPrefix(r[2], "HTTP/")

	return len(req), &RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	// TODO: implement this func
	if r.RequestState == Request_Done {
		return 0, fmt.Errorf("error trying to read data in a done state.")
	}

	if r.RequestState != Request_Initialized {
		if r.RequestState != Request_Done {
			return 0, fmt.Errorf("error: unknown RequestState.")
		}
	}

	for r.RequestState == Request_Initialized {
		pos, req, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing request-line: %s %d", err, pos)
		}

		if pos == 0 {
			return 0, nil
		}
		r.RequestState = Request_Done
		r.RequestLine = *req
	}

	return 0, nil
}

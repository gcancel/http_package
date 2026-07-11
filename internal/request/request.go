package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gcancel/http_package/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       RequestState
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
	Request_Parsing_Headers
)
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, bufferSize)
	request := Request{
		State:   Request_Initialized,
		Headers: headers.NewHeaders(),
	}

	readToIndex := 0

	for request.State != Request_Done {

		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		bytesRead, err := reader.Read(buffer[readToIndex:])

		if errors.Is(err, io.EOF) {
			if request.State != Request_Done {
				return nil, fmt.Errorf("incomplete request: %s", err)
			}
			request.State = Request_Done
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error reading from buffer: %s", err)
		}

		readToIndex += bytesRead

		bytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing bytes: %s Total bytes read: %d Buffer: %d", err, bytesRead, buffer)
		}
		fmt.Println(string(buffer[:bytesParsed]))
		copy(buffer, buffer[bytesParsed:])
		readToIndex -= bytesParsed

	}

	return &request, nil
}

func parseRequestLine(lines []byte) (int, *RequestLine, error) {

	idx := bytes.Index(lines, []byte("\r\n"))
	if idx == -1 {
		return 0, nil, nil
	}

	req := string(lines[:idx])
	r := strings.Split(req, " ")
	if len(r) != 3 {
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

	return idx + 2, &RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalByteParsed := 0
	for r.State != Request_Done {
		n, err := r.parseSingle(data[totalByteParsed:])
		if err != nil {
			return 0, err
		}
		totalByteParsed += n
		if n == 0 {
			break
		}
	}
	return totalByteParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case Request_Initialized:
		bytesParsed, req, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing request-line: %s %d", err, bytesParsed)
		}
		if bytesParsed == 0 {
			return 0, nil
		}

		r.RequestLine = *req
		r.State = Request_Parsing_Headers

		return bytesParsed, nil

	case Request_Parsing_Headers:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing headers: %s", err)
		}
		if done {
			r.State = Request_Done
		}

		return n, nil

	case Request_Done:
		return 0, fmt.Errorf("error trying to read data in a done state.")
	default:
		return 0, fmt.Errorf("error: unknown RequestState.")
	}
}

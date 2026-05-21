package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const (
	crlf = "\r\n"
)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		done = false
		return 0, done, nil
	}
	if idx == 0 {
		done = true
		return n, done, nil
	}

	// turning data into string, splitting it on the space and then trimming the remaining
	headersString := string(data)
	headerLine := strings.Split(headersString, " ")
	if len(headerLine) != 2 {
		done = false
		return 0, done, fmt.Errorf("error: Invalid HTTP field line: %#v", headerLine)
	}

	for _, str := range headerLine {
		str = strings.TrimSpace(str)
	}

	key := headerLine[0]
	key = strings.TrimRight(key, ":")
	value := headerLine[1]
	value = strings.Trim(value, "\r\n")

	h[key] = value
	done = false

	return idx + 2, done, nil
}

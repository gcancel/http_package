package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		done = false
		return 0, done, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headerParts := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(headerParts) != 2 {
		return 0, false, fmt.Errorf("error: Invalid HTTP field line: %#v", string(headerParts[0]))
	}

	key := string(headerParts[0])
	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("error: Invalid header line %s", key)
	}
	key = strings.TrimSpace(key)

	// checking for invalid characters in the field name
	for _, c := range key {
		if !isAlphaNumeric(c) {
			return 0, false, fmt.Errorf("error: Invalid field name characters. %v | %s", c, key)
		}
	}

	key = strings.ToLower(key)
	value := bytes.TrimSpace(headerParts[1])

	h.Set(key, string(value))

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func isAlphaNumeric(c rune) bool {
	if (c >= 'A' && c <= 'z') || (c >= '0' && c <= '9') {
		return true
	}
	return false
}

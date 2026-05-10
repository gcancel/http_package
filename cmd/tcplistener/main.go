package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	// Opening TCP connection...
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("error: tcp connection failed: %s", err.Error())
		}
		fmt.Printf("TCP Connection accepted from %s.\n", conn.RemoteAddr())
		stream := getLinesChannel(conn)

		for line := range stream {
			fmt.Println(line)
		}
		fmt.Printf("TCP Connection closed from %s.\n", conn.RemoteAddr())
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	strChannel := make(chan string)

	go func() {
		defer f.Close()
		defer close(strChannel)
		line := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if n > 0 {
				// process bytes
				// buffer[:n] will store data up to the remaining 'n' left to read!
				line += string(buffer[:n])

				parts := strings.Split(line, "\n")
				if len(parts) > 1 {

					for i := 0; i < len(parts)-1; i++ {
						strChannel <- parts[i]
					}
					line = ""
					line += parts[len(parts)-1]

				}

			}
			if errors.Is(err, io.EOF) {
				break

			}

		}
	}()

	return strChannel
}

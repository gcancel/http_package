package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr := "localhost:42069"
	udp, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	udpConnection, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConnection.Close()

	fmt.Printf("Sending message to %s. Enter Ctrl+C to exit", addr)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = udpConnection.Write([]byte(input))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Message sent: %s\n", input)
	}

}

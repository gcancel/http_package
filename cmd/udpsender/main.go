package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	udpConnection, err := net.DialUDP("udp", nil, udp)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConnection.Close()

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

	}

}

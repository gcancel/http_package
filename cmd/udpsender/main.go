package main

import (
	"log"
	"net"
)

func main() {
	udp, err := net.ResolveUDPAddr("udp", ":")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	port := "8080"

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(fmt.Sprintf("Could not establish TCP server on port %s.", port))
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(fmt.Sprintf("Could not listen on port %s.", port))
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	r := bufio.NewReader(conn)

	s, err := r.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("Error reading from connection: %s", err))
	}
	fmt.Printf("s = %s\n", s)
}

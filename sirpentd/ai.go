package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	//"time"
)

func main() {
	port := os.Args[1] //"8901"

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
	//r := bufio.NewReader(conn)

	for {
		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("SOUTHEAST\n"))
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			//time.Sleep(20000 * time.Millisecond)
			conn.Write([]byte("N"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("O"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("RT"))
			//time.Sleep(2000 * time.Millisecond)
			conn.Write([]byte("H\n"))
			//time.Sleep(2000 * time.Millisecond)
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("NORTHWEST\n"))
		}

		for i := 0; i < 3; i++ {
			bufio.NewReader(conn).ReadString('\n')
			conn.Write([]byte("SOUTH\n"))
		}
	}

	/*s, err := r.ReadString('\n')
	  if err != nil {
	    panic(fmt.Sprintf("Error reading from connection: %s", err))
	  }
	  fmt.Printf("s = %s\n", s)*/
}

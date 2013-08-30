package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8189")
	if err != nil {
		// handle error
	}

	fmt.Println("listen2 fmt...")
	log.Println("listen2...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			// handle error
			continue
		}

		log.Print(conn)
		//go handleConnection(conn)
	}
}

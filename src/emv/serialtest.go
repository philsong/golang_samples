package main

import (
	"github.com/goserial"
	"log"
)

func main() {
	c := &goserial.Config{Name: "COM5", Baud: 115200}
	s, err := goserial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d/%d", n, len(buf))
	log.Printf("%s", buf[:n])
}

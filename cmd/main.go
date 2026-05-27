package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Print("Olaio")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
}

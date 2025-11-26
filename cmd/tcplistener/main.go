package main

import (
	"fmt"
	"log"
	"net"

	"github.com/red7-c/httpfromtcp/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer l.Close()
	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal("error", err)
		}
		req, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatal("error", err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("- Headers:\n")

		req.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
	}
}

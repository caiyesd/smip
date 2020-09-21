package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "p", 60080, "http port")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `smip: smip/0.0.1
Usage: smip [-p port]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalln("failed to listen at", fmt.Sprintf("0.0.0.0:%d, %s", port, err))
		return
	}

	log.Println("start listening at", listener.Addr().String())
	for {
		if conn, err := listener.Accept(); err == nil {
			addr := conn.RemoteAddr().String()
			log.Println(addr)
			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s\n", len([]byte(addr)), addr)))
			conn.Close()
		}
	}
}

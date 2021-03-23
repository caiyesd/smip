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
	fmt.Fprintf(os.Stderr, `smip: smip/0.0.2
Usage: smip [-p port]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	ch := make(chan int)
	go tcpSmip()
	go udpSmip()
	<-ch
}

func tcpSmip() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalln("failed to listen at", fmt.Sprintf("0.0.0.0:%d, %s", port, err))
		return
	}

	log.Println("start listening at tcp port", listener.Addr().String())
	for {
		if conn, err := listener.Accept(); err == nil {
			addr := conn.RemoteAddr().String()
			log.Println("tcp", addr)
			conn.Write([]byte(fmt.Sprintf("%s\n", addr)))
			conn.Close()
		}
	}
}

func udpSmip() {
	listener, err := net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalln("failed to listen at", fmt.Sprintf("0.0.0.0:%d, %s", port, err))
		return
	}
	log.Println("start listening at udp port", fmt.Sprintf("0.0.0.0:%d", port))
	buffer := make([]byte, 4096)

	for {
		_, addr, err := listener.ReadFrom(buffer)
		if err != nil {
			continue
		}
		log.Println("udp", addr.String())
		listener.WriteTo([]byte(fmt.Sprintf("%s\n", addr.String())), addr)
	}

}

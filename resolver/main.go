package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"syscall"
	"time"

	"github.com/benburkert/dns"
)

func bindConnToDevice(conn interface{}, device string) error {
	ptrVal := reflect.ValueOf(conn)
	val := reflect.Indirect(ptrVal)
	//next line will get you the net.netFD
	fdmember := val.FieldByName("fd")
	val1 := reflect.Indirect(fdmember)
	val1 = val1.FieldByName("pfd")
	netFdPtr := val1.FieldByName("Sysfd")
	fd := int(netFdPtr.Int())
	//fd now has the actual fd for the socket

	return syscall.SetsockoptString(fd, syscall.SOL_SOCKET,
		syscall.SO_BINDTODEVICE, device)
}

func bindConnToAddr(conn interface{}, bindAddr net.IP, port int) error {
	ptrVal := reflect.ValueOf(conn)
	val := reflect.Indirect(ptrVal)
	//next line will get you the net.netFD
	fdmember := val.FieldByName("fd")
	val1 := reflect.Indirect(fdmember)
	val1 = val1.FieldByName("pfd")
	netFdPtr := val1.FieldByName("Sysfd")
	fd := int(netFdPtr.Int())
	//fd now has the actual fd for the socket

	// if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
	// 	log.Fatal(err)
	// }
	//syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)

	lsa := syscall.SockaddrInet4{Port: port}
	copy(lsa.Addr[:], bindAddr.To4())
	return syscall.Bind(fd, &lsa)
}

type transport struct {
	ifName string
}

func (t *transport) DialAddr(ctx context.Context, addr net.Addr) (dns.Conn, error) {
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "udp4", addr.String())
	if err != nil {
		return nil, err
	} else {
		if t.ifName != "" {
			err := bindConnToDevice(conn, t.ifName)
			if err != nil {
				return nil, err
			}
			// ifi, err := net.InterfaceByName(t.ifName)
			// if err != nil {
			// 	return nil, err
			// }
			// addrs, err := ifi.Addrs()
			// if err != nil {
			// 	return nil, err
			// }
			// if len(addrs) > 0 {
			// 	fmt.Println(strings.Split(addrs[0].String(), "/")[0])
			// 	err := bindConnToAddr(conn, net.ParseIP(strings.Split(addrs[0].String(), "/")[0]))
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// }
		}
		return &dns.PacketConn{Conn: conn}, nil
	}
}

var (
	server string
	port   int
	ifName string
)

func init() {
	flag.StringVar(&server, "s", "114.114.114.114", "dns server")
	flag.IntVar(&port, "p", 53, "dns port")
	flag.StringVar(&ifName, "I", "", "interface name")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `v0.0.1
Usage: %s [-s dns-server] [-p dns-port] [-I ifname] <domain>
Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	client := new(dns.Client)
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	client.Transport = &transport{ifName: ifName}
	query := &dns.Query{
		RemoteAddr: &net.UDPAddr{IP: net.ParseIP(server), Port: port},
		Message: &dns.Message{
			Questions: []dns.Question{
				{
					Name:  fmt.Sprintf("%s.", flag.Arg(0)),
					Type:  dns.TypeA,
					Class: dns.ClassIN,
				},
			},
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	msg, err := client.Do(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
	hasAnswer := false
	for _, answer := range msg.Answers {
		if dns.TypeA == answer.Record.Type() {
			a := answer.Record.(*dns.A)
			fmt.Println(a.A)
			hasAnswer = true
		}
	}
	if !hasAnswer {
		os.Exit(3)
	}
}

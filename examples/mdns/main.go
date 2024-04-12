package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/cmol/dns"
)

const (
	srvAddr = "224.0.0.251:5353"
	// srvAddr     = "[ff02::fb]:5353"
	maxDatagram = 1500
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [interface_name]\n", os.Args[0])
		os.Exit(1)
	}

	ifi, err := net.InterfaceByName(os.Args[1])
	if err != nil {
		fmt.Printf("Unable to find interface with name: %s \n", os.Args[1])
		fmt.Printf("Usage: %s [interface_name]\n", os.Args[0])
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", srvAddr)
	if err != nil {
		panic(err.Error())
	}

	l, err := net.ListenMulticastUDP("udp", ifi, addr)
	if err != nil {
		panic(err.Error())
	}
	l.SetReadBuffer(maxDatagram)
	for {
		b := make([]byte, maxDatagram)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			fmt.Println(err.Error())
		}
		buf := bytes.NewBuffer(b)
		message, err := dns.ParseMessage(buf)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Read %d bytes from %s%%%s containing: \n%s\n", n,
			src.IP.String(), src.Zone, prettyPrint(message))
	}
}

func prettyPrint(i any) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

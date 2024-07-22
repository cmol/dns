package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/cmol/dns"
	"golang.org/x/net/ipv6"
)

const (
	// srvAddr = "224.0.0.251:5353"
	srvAddr      = "ff02::fb"
	listenAddr   = "[ff02::fb]:5353"
	maxDatagram  = 1024
	responseAddr = "::1"
	dnsName      = "mdns-test.local"
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
	group := net.ParseIP(srvAddr)
	c, err := net.ListenPacket("udp6", listenAddr)
	if err != nil {
		panic(err.Error())
	}
	defer c.Close()
	p := ipv6.NewPacketConn(c)
	if err := p.JoinGroup(ifi, &net.UDPAddr{IP: group}); err != nil {
		panic(err.Error())
	}
	if err := p.SetControlMessage(ipv6.FlagDst, true); err != nil {
		panic(err.Error())
	}
	p.SetTrafficClass(0x0)
	p.SetHopLimit(1)
	p.SetMulticastHopLimit(1)

	for {
		b := make([]byte, maxDatagram)
		n, rcm, src, err := p.ReadFrom(b)
		if err != nil {
			println(err.Error())
			continue
		}
		fmt.Printf("Read %d bytes from %s\n", n, src.String())
		if rcm.Dst.IsMulticast() {
			if rcm.Dst.Equal(group) {
				// joined group, do something
			} else {
				// unknown group, discard
				continue
			}
		} else {
			continue
		}
		buf := bytes.NewBuffer(b)
		message, err := dns.ParseMessage(buf)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Look for requests for our domain
		if message.QR {
			continue
		}
		sendResponse := false
		for _, question := range message.Questions {
			fmt.Printf("Asking for %s\n", question.Domain)
			sendResponse = sendResponse || question.Domain == dnsName
			if sendResponse {
				break
			}
		}

		if !sendResponse {
			continue
		}

		reply := dns.ReplyTo(message)
		record := dns.Record{
			TTL:   60,
			Class: uint16(dns.IN),
			Type:  dns.AAAA,
			Name:  dnsName,
		}
		record.Data = &dns.IPv6{Addr: netip.MustParseAddr(responseAddr)}
		reply.Answers = append(reply.Answers, record)
		buffer := make([]byte, 1024)
		buffer = buffer[:0] // Clear buffer
		sendBuf := bytes.NewBuffer(buffer)
		if err := reply.Build(sendBuf, dns.NewDomains()); err != nil {
			fmt.Println("unable to build message: " + err.Error())
		}
		fmt.Printf("buffer bytes %+v\n", sendBuf.Bytes())

		dst := &net.UDPAddr{IP: group, Port: 5353}
		wcm := ipv6.ControlMessage{TrafficClass: 0xe0, HopLimit: 1}
		for _, ifi := range []*net.Interface{ifi} {
			wcm.IfIndex = ifi.Index
			fmt.Printf("Sending to group %s from interface index %d with data\n%s\n",
				group.String(), ifi.Index, prettyPrint(reply))
			if _, err := p.WriteTo(sendBuf.Bytes(), &wcm, dst); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func prettyPrint(i any) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

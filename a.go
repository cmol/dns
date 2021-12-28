package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"

	"inet.af/netaddr"
)

type IPv4 struct {
	IP netaddr.IP
}

func (ip *IPv4) Parse(buf *bytes.Buffer, domains map[int]string) error {
	a := make([]byte, 4)
	if err := binary.Read(buf, binary.BigEndian, &a); err != nil {
		return errors.New("Could not read IPv4 address")
	}
	ip.IP = netaddr.IPv4(a[0], a[1], a[2], a[3])
	return nil
}

func (ip *IPv4) Build(buf *bytes.Buffer, domains map[string]int) error {
	addr := ip.IP.As4()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

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

func (ip *IPv4) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	var a [4]byte
	if err := binary.Read(buf, binary.BigEndian, a[:]); err != nil {
		return errors.New("Could not read IPv4 address")
	}
	ip.IP = netaddr.IPFrom4(a)
	return nil
}

func (ip *IPv4) Build(buf *bytes.Buffer) error {
	addr := ip.IP.As4()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

func (ip *IPv4) PreBuild(domains *Domains) (int, error) {
	return 4, nil
}

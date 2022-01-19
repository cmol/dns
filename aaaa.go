package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"

	"inet.af/netaddr"
)

type IPv6 struct {
	IP netaddr.IP
}

func (ip *IPv6) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
	var a [16]byte
	if err := binary.Read(buf, binary.BigEndian, a[:]); err != nil {
		return errors.New("Could not read IPv6 address")
	}
	ip.IP = netaddr.IPFrom16(a)
	return nil
}

func (ip *IPv6) Build(buf *bytes.Buffer, domains *Domains) error {
	addr := ip.IP.As16()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

func (ip *IPv6) PreBuild(r *Record, domains *Domains) (int, error) {
	return 16, nil
}

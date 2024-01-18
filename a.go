package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/netip"
)

type IPv4 struct {
	netip.Addr
}

func (ip *IPv4) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
	var a [4]byte
	if err := binary.Read(buf, binary.BigEndian, a[:]); err != nil {
		return errors.New("could not read IPv4 address")
	}
	ip.Addr = netip.AddrFrom4(a)
	return nil
}

func (ip *IPv4) Build(buf *bytes.Buffer, _ *Domains) error {
	addr := ip.As4()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

func (ip *IPv4) PreBuild(_ *Record, _ *Domains) (int, error) {
	return 4, nil
}

package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/netip"
)

// IPv6 implements interface RData
type IPv6 struct {
	netip.Addr
}

// Parse implements AAAA parsing for interface RData
func (ip *IPv6) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
	var a [16]byte
	if err := binary.Read(buf, binary.BigEndian, a[:]); err != nil {
		return errors.New("could not read IPv6 address")
	}
	ip.Addr = netip.AddrFrom16(a)
	return nil
}

// Build implements AAAA building for interface RData
func (ip *IPv6) Build(buf *bytes.Buffer, _ *Domains) error {
	addr := ip.As16()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

// PreBuild step, just returning record size
func (ip *IPv6) PreBuild(_ *Record, _ *Domains) (int, error) {
	return 16, nil
}

// TransformName satisfies the interface
func (*IPv6) TransformName(name string) string { return name }

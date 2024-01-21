package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/netip"
)

// IPv4 implements interface RData
type IPv4 struct {
	netip.Addr
}

// Parse implements A parsing for interface RData
func (ip *IPv4) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
	var a [4]byte
	if err := binary.Read(buf, binary.BigEndian, a[:]); err != nil {
		return errors.New("could not read IPv4 address")
	}
	ip.Addr = netip.AddrFrom4(a)
	return nil
}

// Build implements A building for interface RData
func (ip *IPv4) Build(buf *bytes.Buffer, _ *Domains) error {
	addr := ip.As4()
	binary.Write(buf, binary.BigEndian, addr)
	return nil
}

// PreBuild step, just returning record size
func (ip *IPv4) PreBuild(_ *Record, _ *Domains) (int, error) {
	return 4, nil
}

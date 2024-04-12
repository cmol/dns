package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Opt implements interface for RDATA
type Opt struct {
	UDPSize     uint16
	RCode       byte
	EDNSVersion byte
	DNSSec      bool
	Record      *Record
	Options     map[uint16][]byte
}

// Parse implements OPT parsing for interface RData
func (o *Opt) Parse(buf *bytes.Buffer, _ int, _ *Domains) error {
	o.UDPSize = o.Record.Class
	o.RCode = byte((o.Record.TTL >> 24) & 0xff)
	o.EDNSVersion = byte((o.Record.TTL >> 16) & 0xff)
	if ((o.Record.TTL >> 15) & 0x01) == 0x01 {
		o.DNSSec = true
	}

	readLen := o.Record.Length
	o.Options = map[uint16][]byte{}
	for readLen > 0 {
		var code uint16
		var length uint16
		if err := binary.Read(buf, binary.BigEndian, &code); err != nil {
			return fmt.Errorf("unable to read variable OPT code: %w", err)
		}
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return fmt.Errorf("unable to read variable OPT length: %w", err)
		}
		o.Options[code] = buf.Next(int(length))
		readLen = readLen - length - 4
	}
	return nil
}

// PreBuild implements OPT pre building for interface RData
func (o *Opt) PreBuild(r *Record, _ *Domains) (int, error) {
	var DNSSec uint32
	if o.DNSSec {
		DNSSec = 1
	}
	r.Class = o.UDPSize
	r.TTL = (uint32(o.RCode) << 24) |
		((uint32(o.EDNSVersion) & 0xff) << 16) | (DNSSec << 15)
	return 0, nil
}

// Build implements OPT building for interface RData
func (o *Opt) Build(_ *bytes.Buffer, _ *Domains) error {
	return nil
}

// DefaultOpt returns a standard OPT record
func DefaultOpt(size int) *Record {
	r := &Record{
		Name: "",
		Type: OPT,
	}
	r.Data = &Opt{
		UDPSize: uint16(size),
		Record:  r,
	}
	return r
}

// TransformName satisfies the interface
func (*Opt) TransformName(name string) string { return name }

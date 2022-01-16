package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Opt struct {
	UDPSize     uint16
	RCode       byte
	EDNSVersion byte
	DNSSec      bool
	Record      *Record
	Options     map[uint16][]byte
}

func (o *Opt) Parse(buf *bytes.Buffer, ptr int, domains *Domains) error {
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

func (o *Opt) PreBuild(domains *Domains) (int, error) {
	var DNSSec uint32
	if o.DNSSec {
		DNSSec = 1
	}
	o.Record.Class = o.UDPSize
	o.Record.TTL = (uint32(o.RCode) << 24) |
		((uint32(o.EDNSVersion) & 0xff) << 16) | (DNSSec << 15)
	return 1, nil
}

func (o *Opt) Build(buf *bytes.Buffer, domains *Domains) error {
	return nil
}

package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// CacheFlushBit holds the bit for the mDNS cache flush instruction
const CacheFlushBit = 0x8000

// RData interface for all record types
type RData interface {
	Parse(*bytes.Buffer, int, *Domains) error
	Build(*bytes.Buffer, *Domains) error
	PreBuild(*Record, *Domains) (int, error)
	TransformName(string) string
}

// Record struct used by record specific types
type Record struct {
	TTL        uint32
	Class      uint16
	Length     uint16
	Type       Type
	Name       string
	Data       RData
	CacheFlush bool
}

// ParseRecord is the generic entry to parsing all records
func ParseRecord(buf *bytes.Buffer, ptr int, domains *Domains) (Record, error) {
	var err error
	initialLen := buf.Len()
	r := Record{}
	r.Name, err = ParseName(buf, ptr, domains)
	if err != nil {
		return Record{}, err
	}

	if err = binary.Read(buf, binary.BigEndian, &r.Type); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.Class); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.TTL); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.Length); err != nil {
		return Record{}, err
	}

	// Read and filter out mDNS CacheFlush bit
	r.CacheFlush = (r.Class & CacheFlushBit) == CacheFlushBit
	r.Class = r.Class & 0x7fff

	if err = r.parseRData(buf, ptr+(initialLen-buf.Len()), domains); err != nil {
		return Record{}, err
	}

	return r, nil
}

func (r *Record) parseRData(buf *bytes.Buffer, ptr int, domains *Domains) error {
	var rdata RData
	switch r.Type {
	case A:
		rdata = &IPv4{}
	case AAAA:
		rdata = &IPv6{}
	case OPT:
		rdata = &Opt{Record: r}
	case CNAME:
		rdata = &CName{}
	case PTR:
		rdata = &Ptr{}
	case SRV:
		rdata = &Srv{NameBytes: r.Name}
	default:
		return errors.New("type not supported: " + RRTypeStrings[r.Type])
	}
	err := rdata.Parse(buf, ptr, domains)
	r.Data = rdata
	return err
}

// Build is the generic entry to building all records
func (r *Record) Build(buf *bytes.Buffer, domains *Domains) error {
	r.Name = r.Data.TransformName(r.Name)
	name := BuildName(r.Name, domains)
	domains.SetBuild(buf.Len(), r.Name)
	buf.WriteString(name)
	length, err := r.Data.PreBuild(r, domains)
	if err != nil {
		return err
	}
	r.Length = uint16(length)
	if err := binary.Write(buf, binary.BigEndian, r.Type); err != nil {
		return err
	}
	if r.CacheFlush {
		err = binary.Write(buf, binary.BigEndian, r.Class|CacheFlushBit)
	} else {
		err = binary.Write(buf, binary.BigEndian, r.Class)
	}
	if err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, r.TTL); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, r.Length); err != nil {
		return err
	}
	err = r.Data.Build(buf, domains)
	return err
}

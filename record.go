package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RData interface {
	Parse(*bytes.Buffer, int, *Domains) error
	Build(*bytes.Buffer, *Domains) error
	PreBuild(*Record, *Domains) (int, error)
}

type Record struct {
	TTL    uint32
	Class  uint16
	Length uint16
	Type   Type
	Name   string
	Data   RData
}

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
	default:
		return errors.New("type not supported: " + RRTypeStrings[r.Type])
	}
	err := rdata.Parse(buf, ptr, domains)
	r.Data = rdata
	return err
}

func (r *Record) Build(buf *bytes.Buffer, domains *Domains) error {
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
	if err := binary.Write(buf, binary.BigEndian, r.Class); err != nil {
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

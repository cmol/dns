package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RData interface {
	Parse(*bytes.Buffer, int, *Domains) error
	Build(*bytes.Buffer) error
	PreBuild(*Domains) (int, error)
}

type Record struct {
	TTL         uint32
	Class       uint16
	RDataLength uint16
	RType       Type
	Name        string
	Data        RData
}

func ParseRecord(buf *bytes.Buffer, ptr int, domains *Domains) (Record, error) {
	var err error
	initialLen := buf.Len()
	r := Record{}
	r.Name, err = ParseName(buf, ptr, domains)
	if err != nil {
		return Record{}, err
	}

	if err = binary.Read(buf, binary.BigEndian, &r.RType); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.Class); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.TTL); err != nil {
		return Record{}, err
	}
	if err = binary.Read(buf, binary.BigEndian, &r.RDataLength); err != nil {
		return Record{}, err
	}

	if err = r.parseRData(buf, ptr+(initialLen-buf.Len()), domains); err != nil {
		return Record{}, err
	}

	return r, nil
}

func (r *Record) parseRData(buf *bytes.Buffer, ptr int, domains *Domains) error {
	var rdata RData
	switch r.RType {
	case A:
		rdata = &IPv4{}
	case AAAA:
		rdata = &IPv6{}
	case OPT:
		rdata = &Opt{Record: r}
	default:
		return errors.New("type not supported: " + RRTypeStrings[r.RType])
	}
	err := rdata.Parse(buf, ptr, domains)
	r.Data = rdata
	return err
}

func (r *Record) Build(buf *bytes.Buffer, domains *Domains) error {
	BuildName(buf, r.Name, domains)
	length, err := r.Data.PreBuild(domains)
	if err != nil {
		return err
	}
	r.RDataLength = uint16(length)
	if err := binary.Write(buf, binary.BigEndian, r.RType); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, r.Class); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, r.TTL); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.BigEndian, r.RDataLength); err != nil {
		return err
	}
	if err := r.Data.Build(buf); err != nil {
		return err
	}
	return nil
}

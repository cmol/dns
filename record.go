package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RData interface {
	Parse(*bytes.Buffer, *Domains) error
	Build(*bytes.Buffer, *Domains) error
	PreBuild() error
}

type Record struct {
	TTL         uint32
	Class       uint16
	RDataLength uint16
	RType       Type
	Name        string
	Data        RData
}

func ParseRecord(buf *bytes.Buffer, pointer int, domains *Domains) (Record, error) {
	var err error
	r := Record{}
	r.Name, err = ParseName(buf, pointer, domains)
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

	if err = r.parseRData(buf, domains); err != nil {
		return Record{}, err
	}

	return r, nil
}

func (r *Record) parseRData(buf *bytes.Buffer, domains *Domains) error {
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
	err := rdata.Parse(buf, domains)
	r.Data = rdata
	return err
}

func (r *Record) Build(buf *bytes.Buffer, domains *Domains) error {
	BuildName(buf, r.Name, domains)
	if err := r.Data.PreBuild(); err != nil {
		return err
	}
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
	if err := r.Data.Build(buf, domains); err != nil {
		return err
	}
	return nil
}

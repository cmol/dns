package dnsmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RData interface {
	Parse(*bytes.Buffer, map[int]string) error
	Build(*bytes.Buffer, map[string]int) error
}

type Record struct {
	TTL         uint32
	Class       uint16
	RDataLength uint16
	RType       Type
	Name        string
	Data        RData
}

func ParseRecord(buf *bytes.Buffer, pointer int, domains map[int]string) (Record, error) {
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

	if r.Data, err = parseRData(buf, r.RType, domains); err != nil {
		return Record{}, err
	}

	return r, nil
}

func parseRData(buf *bytes.Buffer, typ Type, domains map[int]string) (RData, error) {
	var rdata RData
	switch typ {
	case A:
		rdata = &IPv4{}
	case AAAA:
		rdata = &IPv6{}
	default:
		return rdata, errors.New("Type not supported")
	}
	err := rdata.Parse(buf, domains)
	return rdata, err
}

func (r *Record) Build(buf *bytes.Buffer, domains map[string]int) error {
	BuildName(buf, r.Name, domains)
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

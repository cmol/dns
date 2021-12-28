package dnsmessage

import (
	"bytes"
	"encoding/binary"
)

const HDR_LENGTH = 12

type Message struct {
	id          uint16
	qr          bool
	opcode      uint8
	aa          bool
	tc          bool
	rd          bool
	ra          bool
	rcode       uint8
	qdcount     uint16
	ancount     uint16
	nscount     uint16
	arcount     uint16
	questions   []Question
	answers     []Record
	nameservers []Record
	additional  []Record
}

func (m *Message) ParseHeader(buf *bytes.Buffer) error {
	var opts uint16
	var err error
	err = binary.Read(buf, binary.BigEndian, &m.id)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &opts)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &m.qdcount)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &m.ancount)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &m.nscount)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &m.arcount)
	if err != nil {
		return err
	}
	m.parseOpts(opts)
	return nil
}

func (m *Message) BuildHeader(buf *bytes.Buffer) error {
	var err error
	err = binary.Write(buf, binary.BigEndian, m.id)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, m.buildOpts())
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, m.qdcount)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, m.ancount)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, m.nscount)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, m.arcount)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) parseOpts(opts uint16) {
	m.rcode = uint8(opts & 0xf)
	m.opcode = uint8((opts >> 11) & 0xf)
	m.ra = opts&OPT_RA == OPT_RA
	m.rd = opts&OPT_RD == OPT_RD
	m.tc = opts&OPT_TC == OPT_TC
	m.aa = opts&OPT_AA == OPT_AA
	m.qr = opts&OPT_QR == OPT_QR
}

func (m *Message) buildOpts() uint16 {
	var opts uint16
	opts |= (uint16(m.rcode) & 0x000f)
	opts |= (uint16(m.opcode) & 0x000f) << 11
	opts |= opt(m.ra, OPT_RA)
	opts |= opt(m.rd, OPT_RD)
	opts |= opt(m.tc, OPT_TC)
	opts |= opt(m.aa, OPT_AA)
	opts |= opt(m.qr, OPT_QR)
	return opts
}

func opt(flag bool, bit uint16) uint16 {
	if flag {
		return bit
	}
	return 0
}

func ParseMessage(buf *bytes.Buffer) (*Message, error) {
	m := new(Message)
	var err error
	err = m.ParseHeader(buf)
	if err != nil {
		return m, err
	}
	ptr := HDR_LENGTH
	//bufLen := buf.Len()
	domains := map[int]string{}
	err = m.parseQuestions(buf, domains, ptr)
	if err != nil {
		return m, err
	}
	err = m.parseRecords(buf, domains, ptr, m.ancount, m.answers)
	if err != nil {
		return m, err
	}
	err = m.parseRecords(buf, domains, ptr, m.nscount, m.nameservers)
	if err != nil {
		return m, err
	}
	err = m.parseRecords(buf, domains, ptr, m.arcount, m.additional)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (m *Message) parseQuestions(buf *bytes.Buffer, domains map[int]string, ptr int) error {
	bufLen := buf.Len()
	for i := 0; i < int(m.qdcount); i++ {
		q, err := ParseQuestion(buf, ptr, domains)
		if err != nil {
			return err
		}
		m.questions = append(m.questions, q)
		ptr = ptr + bufLen - buf.Len()
		bufLen = buf.Len()
	}
	return nil
}

func (m *Message) parseRecords(buf *bytes.Buffer, domains map[int]string, ptr int,
	count uint16, list []Record) error {
	bufLen := buf.Len()
	for i := 0; i < int(count); i++ {
		r, err := ParseRecord(buf, ptr, domains)
		if err != nil {
			return err
		}
		list = append(list, r)
		ptr = ptr + bufLen - buf.Len()
		bufLen = buf.Len()
	}
	return nil
}

package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const HdrLength = 12

type Message struct {
	ID          uint16
	QR          bool
	OPCode      uint8
	AA          bool
	TC          bool
	RD          bool
	RA          bool
	RCode       uint8
	qdcount     uint16
	ancount     uint16
	nscount     uint16
	arcount     uint16
	Questions   []Question
	Answers     []Record
	Nameservers []Record
	Additional  []Record
}

func (m *Message) ParseHeader(buf *bytes.Buffer) error {
	var opts uint16
	var err error
	err = binary.Read(buf, binary.BigEndian, &m.ID)
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
	err = binary.Write(buf, binary.BigEndian, m.ID)
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
	m.RCode = uint8(opts & 0xf)
	m.OPCode = uint8((opts >> 11) & 0xf)
	m.RA = opts&OptRa == OptRa
	m.RD = opts&OptRd == OptRd
	m.TC = opts&OptTc == OptTc
	m.AA = opts&OptAa == OptAa
	m.QR = opts&OptQr == OptQr
}

func (m *Message) buildOpts() uint16 {
	var opts uint16
	opts |= (uint16(m.RCode) & 0x000f)
	opts |= (uint16(m.OPCode) & 0x000f) << 11
	opts |= opt(m.RA, OptRa)
	opts |= opt(m.RD, OptRd)
	opts |= opt(m.TC, OptTc)
	opts |= opt(m.AA, OptAa)
	opts |= opt(m.QR, OptQr)
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
	ptr := HdrLength
	domains := &Domains{parsePtr: map[int]string{}, buildPtr: map[string]int{}}
	err = m.parseQuestions(buf, domains, ptr)
	if err != nil {
		return m, fmt.Errorf("unable to parse questions: %w", err)
	}
	m.Answers, err = m.parseRecords(buf, domains, ptr, m.ancount, m.Answers)
	if err != nil {
		return m, fmt.Errorf("unable to parse answers: %w", err)
	}
	m.Nameservers, err = m.parseRecords(buf, domains, ptr, m.nscount, m.Nameservers)
	if err != nil {
		return m, fmt.Errorf("unable to parse nameservers: %w", err)
	}
	m.Additional, err = m.parseRecords(buf, domains, ptr, m.arcount, m.Additional)
	if err != nil {
		return m, fmt.Errorf("unable to parse additionals: %w", err)
	}
	return m, nil
}

func (m *Message) parseQuestions(buf *bytes.Buffer, domains *Domains, ptr int) error {
	bufLen := buf.Len()
	for i := 0; i < int(m.qdcount); i++ {
		q, err := ParseQuestion(buf, ptr, domains)
		if err != nil {
			return err
		}
		m.Questions = append(m.Questions, q)
		ptr = ptr + bufLen - buf.Len()
		bufLen = buf.Len()
	}
	return nil
}

func (m *Message) parseRecords(buf *bytes.Buffer, domains *Domains, ptr int,
	count uint16, list []Record) ([]Record, error) {
	bufLen := buf.Len()
	for i := 0; i < int(count); i++ {
		r, err := ParseRecord(buf, ptr, domains)
		if err != nil {
			return list, err
		}
		list = append(list, r)
		ptr = ptr + bufLen - buf.Len()
		bufLen = buf.Len()
	}
	return list, nil
}

func (m *Message) Build(buf *bytes.Buffer, domains *Domains) error {
	m.qdcount = uint16(len(m.Questions))
	m.ancount = uint16(len(m.Answers))
	m.nscount = uint16(len(m.Nameservers))
	m.arcount = uint16(len(m.Additional))
	err := m.BuildHeader(buf)
	if err != nil {
		return errors.New("Unable to build header")
	}

	err = m.buildQuestions(buf, domains)
	if err != nil {
		return err
	}
	err = m.buildRecords(buf, domains, m.Answers)
	if err != nil {
		return err
	}
	err = m.buildRecords(buf, domains, m.Nameservers)
	if err != nil {
		return err
	}
	err = m.buildRecords(buf, domains, m.Additional)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) buildQuestions(buf *bytes.Buffer, domains *Domains) error {
	for _, q := range m.Questions {
		err := q.Build(buf, domains)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) buildRecords(buf *bytes.Buffer, domains *Domains, records []Record) error {
	for _, r := range records {
		err := r.Build(buf, domains)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReplyTo(other *Message) *Message {
	m := &Message{
		ID:        other.ID,
		QR:        true,
		Questions: other.Questions,
	}
	return m
}
